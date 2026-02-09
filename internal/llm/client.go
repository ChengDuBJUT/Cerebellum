package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaClient handles communication with Ollama
type OllamaClient struct {
	host   string
	model  string
	client *http.Client
}

// NewOllama creates a new Ollama client
func NewOllama(host, model string) *OllamaClient {
	return &OllamaClient{
		host:   host,
		model:  model,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// Generate sends a prompt to Ollama and returns the response
func (c *OllamaClient) Generate(prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model":  c.model,
		"prompt": prompt,
		"stream": false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.host+"/api/generate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if response, ok := result["response"].(string); ok {
		return response, nil
	}

	return "", fmt.Errorf("no response in result")
}

// GenerateStream sends a prompt and returns a channel of response chunks
func (c *OllamaClient) GenerateStream(prompt string) (<-chan string, error) {
	reqBody := map[string]interface{}{
		"model":  c.model,
		"prompt": prompt,
		"stream": true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.host+"/api/generate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	ch := make(chan string)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		decoder := json.NewDecoder(resp.Body)
		for {
			var result map[string]interface{}
			if err := decoder.Decode(&result); err != nil {
				break
			}
			if response, ok := result["response"].(string); ok {
				ch <- response
			}
			if done, ok := result["done"].(bool); ok && done {
				break
			}
		}
	}()

	return ch, nil
}

// GetModel returns the model name
func (c *OllamaClient) GetModel() string {
	return c.model
}

// GetHost returns the host URL
func (c *OllamaClient) GetHost() string {
	return c.host
}
