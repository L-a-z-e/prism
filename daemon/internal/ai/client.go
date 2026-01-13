package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client interface {
	GenerateCode(ctx context.Context, prompt string) (string, error)
	ReviewCode(ctx context.Context, code string) (string, error)
	GenerateCodeStream(ctx context.Context, prompt string, callback func(chunk string) error) error
}

// ============================================================================
// Ollama Client - 실제 DeepSeek/로컬 LLM 사용
// ============================================================================

type OllamaClient struct {
	baseURL string
	model   string
	timeout time.Duration
	client  *http.Client
}

func NewOllamaClient() *OllamaClient {
	baseURL := os.Getenv("OLLAMA_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "deepseek-r1:14b"
	}

	timeout := 300 * time.Second
	if timeoutEnv := os.Getenv("OLLAMA_TIMEOUT"); timeoutEnv != "" {
		if d, err := time.ParseDuration(timeoutEnv); err == nil {
			timeout = d
		}
	}

	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type OllamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type OllamaResponse struct {
	Model         string `json:"model"`
	Response      string `json:"response"`
	Done          bool   `json:"done"`
	Context       []int  `json:"context,omitempty"`
	TotalDuration int64  `json:"total_duration,omitempty"`
}

// GenerateCode - 비스트리밍 방식 (전체 응답 대기)
func (c *OllamaClient) GenerateCode(ctx context.Context, prompt string) (string, error) {
	var result string
	var lastErr error

	// 재시도 로직 (최대 3회)
	for attempt := 1; attempt <= 3; attempt++ {
		req := OllamaRequest{
			Model:  c.model,
			Prompt: prompt,
			Stream: false,
			Options: map[string]interface{}{
				"temperature": 0.1,
				"top_p":       0.9,
				"top_k":       40,
			},
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			return "", fmt.Errorf("failed to marshal request: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(httpReq)
		if err != nil {
			lastErr = err
			if attempt < 3 {
				// exponential backoff
				backoff := time.Duration(1<<uint(attempt-1)) * time.Second
				select {
				case <-time.After(backoff):
					continue
				case <-ctx.Done():
					return "", ctx.Err()
				}
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
			if attempt < 3 && resp.StatusCode >= 500 {
				backoff := time.Duration(1<<uint(attempt-1)) * time.Second
				select {
				case <-time.After(backoff):
					continue
				case <-ctx.Done():
					return "", ctx.Err()
				}
			}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		var ollamaResp OllamaResponse
		if err := json.Unmarshal(body, &ollamaResp); err != nil {
			lastErr = fmt.Errorf("failed to unmarshal response: %w", err)
			continue
		}

		result = ollamaResp.Response
		return result, nil
	}

	return "", fmt.Errorf("all retry attempts failed: %w", lastErr)
}

// GenerateCodeStream - 스트리밍 방식 (실시간 응답)
func (c *OllamaClient) GenerateCodeStream(ctx context.Context, prompt string, callback func(chunk string) error) error {
	req := OllamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: true, // 스트리밍 활성화
		Options: map[string]interface{}{
			"temperature": 0.1,
			"top_p":       0.9,
			"top_k":       40,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	// 스트리밍 응답 처리
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Bytes()
		var streamResp OllamaResponse
		if err := json.Unmarshal(line, &streamResp); err != nil {
			continue
		}

		if streamResp.Response != "" {
			if err := callback(streamResp.Response); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func (c *OllamaClient) ReviewCode(ctx context.Context, code string) (string, error) {
	prompt := fmt.Sprintf(`Review the following code and provide detailed feedback in Korean:

Code:
%s

Please provide:
1. Code quality assessment
2. Potential bugs or issues
3. Performance improvements
4. Best practice suggestions`, code)

	return c.GenerateCode(ctx, prompt)
}

// ============================================================================
// Mock Client - 테스트용 (OLLAMA_MOCK=true 또는 Ollama 연결 불가능 시)
// ============================================================================

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (c *MockClient) GenerateCode(ctx context.Context, prompt string) (string, error) {
	// 실제 요청 흔적을 보이기 위해 일부 프롬프트 내용 반영
	lang := "Java"
	if strings.Contains(prompt, "TypeScript") || strings.Contains(prompt, "Vue") {
		lang = "TypeScript"
	} else if strings.Contains(prompt, "Python") {
		lang = "Python"
	}

	codeSnippet := fmt.Sprintf(`// Mock generated code (%s)
// Based on prompt: %s

package com.example.demo;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api")
public class DemoController {
    
    @GetMapping("/hello")
    public String hello() {
        return "Hello from Mock AI!";
    }
}`, lang, strings.TrimSpace(prompt)[:50])

	// Mock 실행 시간 시뮬레이션
	select {
	case <-time.After(500 * time.Millisecond):
		return codeSnippet, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (c *MockClient) ReviewCode(ctx context.Context, code string) (string, error) {
	return "✅ Mock Review: Code looks good! No issues found.", nil
}

func (c *MockClient) GenerateCodeStream(ctx context.Context, prompt string, callback func(chunk string) error) error {
	chunks := []string{
		"// Generating code...\n",
		"package com.example;\n\n",
		"public class Demo {\n",
		"    public static void main(String[] args) {\n",
		"        System.out.println(\"Hello\");\n",
		"    }\n",
		"}\n",
	}

	for _, chunk := range chunks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := callback(chunk); err != nil {
			return err
		}

		// 실시간 스트리밍 시뮬레이션
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// ============================================================================
// Client Factory - 환경 변수에 따라 적절한 클라이언트 생성
// ============================================================================

func NewClient() Client {
	// 강제로 Mock 사용
	if os.Getenv("OLLAMA_MOCK") == "true" {
		fmt.Println("📋 Using Mock Client (OLLAMA_MOCK=true)")
		return NewMockClient()
	}

	// 실제 Ollama 클라이언트 시도
	ollama := NewOllamaClient()

	// 연결 가능 여부 확인 (헬스 체크)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	healthCheckReq, _ := http.NewRequestWithContext(ctx, "GET", ollama.baseURL+"/api/tags", nil)
	resp, err := ollama.client.Do(healthCheckReq)
	if err == nil && resp.StatusCode == http.StatusOK {
		resp.Body.Close()
		fmt.Printf("✅ Ollama connected at %s (model: %s)\n", ollama.baseURL, ollama.model)
		return ollama
	}

	if resp != nil {
		resp.Body.Close()
	}

	// Ollama 연결 실패 시 Mock으로 폴백
	fmt.Printf("⚠️  Ollama not available at %s, falling back to Mock Client\n", ollama.baseURL)
	return NewMockClient()
}
