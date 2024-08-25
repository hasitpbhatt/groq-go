package groq

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatCompletion(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Mock server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id": "123", "object": "text", "created": 1643723400, "model": "llama3-8b-8192", "choices": [{"index": 0, "message": {"role": "user", "content": "Hello, world!"}, "logprobs": null, "finish_reason": "length"}], "usage": {"queue_time": 0.1, "prompt_tokens": 5, "prompt_time": 0.2, "completion_tokens": 10, "completion_time": 0.3, "total_tokens": 15, "total_time": 0.6}, "system_fingerprint": "1234567890", "x_groq": {"id": "123"}}`))
		}))
		defer ts.Close()

		c := NewClient(ts.Client(), "test-key")
		c.chatCompletionURL = ts.URL

		// Test data
		messages := []Message{
			{Role: "user", Content: "Hello, world!"},
		}

		// Call the function under test
		completion, err := c.ChatCompletion(messages)

		// Assertions
		assert.Nil(t, err)
		assert.NotNil(t, completion)
		assert.Equal(t, "123", completion.ID)
		assert.Equal(t, "text", completion.Object)
		assert.Equal(t, 1643723400, completion.Created)
		assert.Equal(t, "llama3-8b-8192", completion.Model)
		assert.Equal(t, 1, len(completion.Choices))
		assert.Equal(t, "Hello, world!", completion.Choices[0].Message.Content)
		assert.Equal(t, "length", completion.Choices[0].FinishReason)
	})

	t.Run("Error", func(t *testing.T) {
		// Mock server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		// Set API key for the test
		os.Setenv("GROQ_API_KEY", "test-key")

		// Test data
		messages := []Message{
			{Role: "user", Content: "Hello, world!"},
		}

		c := Client{
			httpClient: ts.Client(),
		}

		// Call the function under test
		completion, err := c.ChatCompletion(messages)

		fmt.Println(completion)
		fmt.Println(err)
		// Assertions
		assert.NotNil(t, err)
	})
}
