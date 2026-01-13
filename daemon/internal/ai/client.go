// daemon/internal/ai/client.go
package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client interface {
	GenerateCode(ctx context.Context, prompt string) (string, error)
	GenerateCodeWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error)
	GenerateCodeStream(ctx context.Context, prompt string, callback func(chunk string) error) error
	GenerateCodeStreamWithSystem(ctx context.Context, systemPrompt, userPrompt string, callback func(chunk string) error) error
	ReviewCode(ctx context.Context, code string) (string, error)
}

// ============================================================================
// Ollama Client
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

	timeout := 600 * time.Second
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

// Generate API types
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

// Chat API types
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string                 `json:"model"`
	Messages []ChatMessage          `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type ChatResponse struct {
	Model   string      `json:"model"`
	Message ChatMessage `json:"message"`
	Done    bool        `json:"done"`
}

func (c *OllamaClient) GenerateCode(ctx context.Context, prompt string) (string, error) {
	return c.retryRequest(ctx, func() (string, error) {
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

		return c.callGenerateAPI(ctx, req)
	})
}

func (c *OllamaClient) GenerateCodeWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	return c.retryRequest(ctx, func() (string, error) {
		req := ChatRequest{
			Model: c.model,
			Messages: []ChatMessage{
				{Role: "system", Content: systemPrompt},
				{Role: "user", Content: userPrompt},
			},
			Stream: false,
			Options: map[string]interface{}{
				"temperature": 0.1,
				"top_p":       0.9,
				"top_k":       40,
				"num_predict": 8192,
			},
		}

		return c.callChatAPI(ctx, req)
	})
}

func (c *OllamaClient) GenerateCodeStream(ctx context.Context, prompt string, callback func(chunk string) error) error {
	req := OllamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: true,
		Options: map[string]interface{}{
			"temperature": 0.1,
			"top_p":       0.9,
			"top_k":       40,
		},
	}

	return c.streamGenerateAPI(ctx, req, callback)
}

func (c *OllamaClient) GenerateCodeStreamWithSystem(ctx context.Context, systemPrompt, userPrompt string, callback func(chunk string) error) error {
	req := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: true,
		Options: map[string]interface{}{
			"temperature": 0.1,
			"top_p":       0.9,
			"top_k":       40,
			"num_predict": 8192,
		},
	}

	return c.streamChatAPI(ctx, req, callback)
}

func (c *OllamaClient) ReviewCode(ctx context.Context, code string) (string, error) {
	systemPrompt := "You are an expert code reviewer. Provide concise, actionable feedback in Korean."
	userPrompt := fmt.Sprintf("Review this code:\n\n%s\n\nProvide:\n1. Quality assessment\n2. Bugs or issues\n3. Performance improvements\n4. Best practices", code)

	return c.GenerateCodeWithSystem(ctx, systemPrompt, userPrompt)
}

// retryRequest wraps any request with exponential backoff retry logic
func (c *OllamaClient) retryRequest(ctx context.Context, fn func() (string, error)) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		lastErr = err

		if attempt < 3 {
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-time.After(backoff):
				continue
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}

	return "", fmt.Errorf("all retry attempts failed: %w", lastErr)
}

func (c *OllamaClient) callGenerateAPI(ctx context.Context, req OllamaRequest) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	return ollamaResp.Response, nil
}

func (c *OllamaClient) callChatAPI(ctx context.Context, req ChatRequest) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	return chatResp.Message.Content, nil
}

func (c *OllamaClient) streamGenerateAPI(ctx context.Context, req OllamaRequest, callback func(chunk string) error) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var streamResp OllamaResponse
		if err := json.Unmarshal(scanner.Bytes(), &streamResp); err != nil {
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

func (c *OllamaClient) streamChatAPI(ctx context.Context, req ChatRequest, callback func(chunk string) error) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var streamResp ChatResponse
		if err := json.Unmarshal(scanner.Bytes(), &streamResp); err != nil {
			continue
		}

		if streamResp.Message.Content != "" {
			if err := callback(streamResp.Message.Content); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

// ============================================================================
// Mock Client
// ============================================================================

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (c *MockClient) GenerateCode(ctx context.Context, prompt string) (string, error) {
	return c.mockResponse(ctx, "Generate API", prompt)
}

func (c *MockClient) GenerateCodeWithSystem(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	return c.mockResponse(ctx, "Chat API", systemPrompt+"\n\n"+userPrompt)
}

func (c *MockClient) ReviewCode(ctx context.Context, code string) (string, error) {
	return "✅ Mock Review: Code looks good!", nil
}

func (c *MockClient) GenerateCodeStream(ctx context.Context, prompt string, callback func(chunk string) error) error {
	return c.mockStream(ctx, callback)
}

func (c *MockClient) GenerateCodeStreamWithSystem(ctx context.Context, systemPrompt, userPrompt string, callback func(chunk string) error) error {
	return c.mockStream(ctx, callback)
}

func (c *MockClient) mockResponse(ctx context.Context, apiType, prompt string) (string, error) {
	code := fmt.Sprintf(`package com.example.demo;

import org.springframework.web.bind.annotation.*;
import org.springframework.http.ResponseEntity;

@RestController
@RequestMapping("/api/mock")
public class MockController {
    
    @GetMapping
    public ResponseEntity<String> get() {
        return ResponseEntity.ok("Mock from %s");
    }
    
    @PostMapping
    public ResponseEntity<String> create(@RequestBody String data) {
        return ResponseEntity.ok("Created: " + data);
    }
}`, apiType)

	select {
	case <-time.After(500 * time.Millisecond):
		return code, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (c *MockClient) mockStream(ctx context.Context, callback func(chunk string) error) error {
	chunks := []string{
		"package com.example;\n\n",
		"public class Demo {\n",
		"    public static void main(String[] args) {\n",
		"        System.out.println(\"Mock Stream\");\n",
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
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// ============================================================================
// Factory
// ============================================================================

func NewClient() Client {
	if os.Getenv("OLLAMA_MOCK") == "true" {
		fmt.Println("📋 Using Mock Client (OLLAMA_MOCK=true)")
		return NewMockClient()
	}

	ollama := NewOllamaClient()

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

	fmt.Printf("⚠️  Ollama not available at %s, falling back to Mock\n", ollama.baseURL)
	return NewMockClient()
}
