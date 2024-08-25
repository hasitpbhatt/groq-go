package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// NewClient creates a new client for interacting with the Groq API.
// It takes the API key as a parameter and returns a pointer to the client.
func NewClient(apiKey string) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("GROQ_API_KEY")
	}
	return &Client{
		apiKey:            apiKey,
		httpClient:        &http.Client{}, // Initialize the HTTP client
		chatCompletionURL: "https://api.groq.com/openai/v1/chat/completions",
	}
}

// ChatCompletion is a function that sends a request to the Groq API for chat completions.
// It takes a slice of Message as input and returns a pointer to http.Response and an error.
func (c *Client) ChatCompletion(messages []Message, options ...Option) (*ChatCompletionResponse, error) {
	body := requestBody{
		Messages:    messages,
		Model:       "llama3-8b-8192",
		Temperature: 1,
		MaxTokens:   1024,
		TopP:        1,
		Stream:      false,
		Stop:        nil,
	}

	for _, option := range options {
		option(&body)
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.chatCompletionURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	completion := ChatCompletionResponse{}
	err = json.NewDecoder(resp.Body).Decode(&completion)
	if err != nil {
		return nil, err
	}

	return &completion, nil
}

// WithModel sets the model for the request body.
func WithModel(model string) func(*requestBody) {
	return func(rb *requestBody) {
		rb.Model = model
	}
}

// WithTemperature sets the temperature for the request body.
func WithTemperature(temperature float64) func(*requestBody) {
	return func(rb *requestBody) {
		rb.Temperature = temperature
	}
}

// WithMaxTokens sets the maximum number of tokens for the request body.
func WithMaxTokens(maxTokens int) func(*requestBody) {
	return func(rb *requestBody) {
		rb.MaxTokens = maxTokens
	}
}

// WithTopP sets the top_p value for the request body.
func WithTopP(topP float64) func(*requestBody) {
	return func(rb *requestBody) {
		rb.TopP = topP
	}
}
