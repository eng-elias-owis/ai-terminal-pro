package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents the LiteLLM API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new LiteLLM API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CompletionRequest represents an API request
type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResponse represents an API response
type CompletionResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Context contains terminal context for AI
 type Context struct {
	OS          string
	Shell       string
	WorkingDir  string
}

// GenerateCommand creates an AI-generated command from natural language
func (c *Client) GenerateCommand(ctx context.Context, userPrompt string, context Context) (string, error) {
	// Build system prompt with context
	systemPrompt := fmt.Sprintf(`You are a terminal command generator. 
Generate the correct command for the user's request.

Context:
- OS: %s
- Shell: %s
- Current Directory: %s

Rules:
1. Return ONLY the command, no explanations
2. Use appropriate syntax for the detected shell
3. Ensure paths are properly escaped
4. Use OS-appropriate commands (e.g., 'dir' for Windows CMD, 'ls' for bash)

Generate command:`, context.OS, context.Shell, context.WorkingDir)

	req := CompletionRequest{
		Model:       "qwen3-terminal",
		MaxTokens:   100,
		Temperature: 0.1,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	command := resp.Choices[0].Message.Content
	// Clean up the command (remove markdown, extra whitespace)
	command = cleanCommand(command)

	return command, nil
}

// sendRequest sends the API request to LiteLLM
func (c *Client) sendRequest(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", httpResp.StatusCode, string(body))
	}

	var resp CompletionResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

// cleanCommand removes markdown and whitespace from AI output
func cleanCommand(cmd string) string {
	// Remove code blocks
	cmd = strings.TrimPrefix(cmd, "```bash")
	cmd = strings.TrimPrefix(cmd, "```sh")
	cmd = strings.TrimPrefix(cmd, "```")
	cmd = strings.TrimSuffix(cmd, "```")
	
	// Trim whitespace
	cmd = strings.TrimSpace(cmd)
	
	return cmd
}
