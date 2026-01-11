package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client interface {
	GenerateCode(prompt string) (string, error)
	ReviewCode(code string) (string, error)
}

// Ollama Client
type OllamaClient struct {
	baseURL string
	model   string
	timeout time.Duration
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

	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		timeout: 300 * time.Second,
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

func (c *OllamaClient) GenerateCode(prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1, // 코드 생성은 일관성 중요
			"top_p":       0.9,
			"top_k":       40,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{
		Timeout: c.timeout,
	}

	resp, err := client.Post(
		c.baseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to call ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return ollamaResp.Response, nil
}

func (c *OllamaClient) ReviewCode(code string) (string, error) {
	prompt := fmt.Sprintf(`Review the following code and provide detailed feedback:

Code:
%s

Please provide:
1. Code quality assessment
2. Potential bugs or issues
3. Performance improvements
4. Best practice suggestions`, code)

	return c.GenerateCode(prompt)
}

// Mock Client (테스트용)
type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (c *MockClient) GenerateCode(prompt string) (string, error) {
	return `// Mock generated code
package com.example.demo;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api")
public class DemoController {
    
    @GetMapping("/hello")
    public String hello() {
        return "Hello from Mock AI!";
    }
}`, nil
}

func (c *MockClient) ReviewCode(code string) (string, error) {
	return "✅ Mock Review: Code looks good! No issues found.", nil
}
