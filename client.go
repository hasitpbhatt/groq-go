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
func NewClient(httpClient *http.Client, apiKey string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{} // Use default client if none provided
	}
	if apiKey == "" {
		apiKey = os.Getenv("GROQ_API_KEY")
	}
	return &Client{
		apiKey:            apiKey,
		httpClient:        httpClient, // Initialize the HTTP client
		chatCompletionURL: "https://api.groq.com/openai/v1/chat/completions",
	}
}

// Client represents a client for interacting with the Groq API.
type Client struct {
	apiKey            string
	httpClient        *http.Client // Added field for HTTP client
	chatCompletionURL string       // Added field for chat completion URL
}

// SetAPIKey sets the API key for the client.
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

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

// ChatCompletionResponse represents the structure of the response received from the Groq API for chat completions.
// It contains the ID of the completion, the object type, the creation time, the model used, the choices made, the usage statistics, the system fingerprint, and the x_groq information.
type ChatCompletionResponse struct {
	ID      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int    `json:"created,omitempty"`
	Model   string `json:"model,omitempty"`
	Choices []struct {
		Index        int         `json:"index,omitempty"`
		Message      Message     `json:"message,omitempty"`
		Logprobs     interface{} `json:"logprobs,omitempty"`
		FinishReason string      `json:"finish_reason,omitempty"`
	} `json:"choices,omitempty"`
	Usage struct {
		QueueTime        float64 `json:"queue_time,omitempty"`
		PromptTokens     int     `json:"prompt_tokens,omitempty"`
		PromptTime       float64 `json:"prompt_time,omitempty"`
		CompletionTokens int     `json:"completion_tokens,omitempty"`
		CompletionTime   float64 `json:"completion_time,omitempty"`
		TotalTokens      int     `json:"total_tokens,omitempty"`
		TotalTime        float64 `json:"total_time,omitempty"`
	} `json:"usage,omitempty"`
	SystemFingerprint string `json:"system_fingerprint,omitempty"`
	XGroq             struct {
		ID string `json:"id,omitempty"`
	} `json:"x_groq,omitempty"`
}

// ChatCompletion is a function that sends a request to the Groq API for chat completions.
// It takes a slice of Message as input and returns a pointer to http.Response and an error.
func (c *Client) ChatCompletion(messages []Message, options ...func(*RequestBody)) (*ChatCompletionResponse, error) {
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
