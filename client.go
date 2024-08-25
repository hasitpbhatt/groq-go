package groq

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

// Message represents a single message in the chat completion request.
// It contains the role of the message sender (e.g., user or system) and the content of the message itself.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RequestBody represents the structure of the request body sent to the Groq API for chat completion.
type RequestBody struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	TopP        float64   `json:"top_p"`
	Stream      bool      `json:"stream"`
	Stop        *string   `json:"stop,omitempty"`
}

// ChatCompletion is a function that sends a request to the Groq API for chat completions.
// It takes a slice of Message as input and returns a pointer to http.Response and an error.
func ChatCompletion(messages []Message, options ...func(*RequestBody)) (*http.Response, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"
	apiKey := os.Getenv("GROQ_API_KEY")

	body := RequestBody{
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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	return client.Do(req)
}

// WithModel sets the model for the request body.
func WithModel(model string) func(*RequestBody) {
	return func(rb *RequestBody) {
		rb.Model = model
	}
}

// WithTemperature sets the temperature for the request body.
func WithTemperature(temperature float64) func(*RequestBody) {
	return func(rb *RequestBody) {
		rb.Temperature = temperature
	}
}

// WithMaxTokens sets the maximum number of tokens for the request body.
func WithMaxTokens(maxTokens int) func(*RequestBody) {
	return func(rb *RequestBody) {
		rb.MaxTokens = maxTokens
	}
}

// WithTopP sets the top_p value for the request body.
func WithTopP(topP float64) func(*RequestBody) {
	return func(rb *RequestBody) {
		rb.TopP = topP
	}
}
