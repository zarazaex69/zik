package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zarazaex69/zik/apps/server/internal/api"
	"github.com/zarazaex69/zik/apps/server/internal/config"
	"github.com/zarazaex69/zik/apps/server/internal/domain"
)

// --- Mock Implementations ---

// MockAIClient is a mock implementation of ai.AIClienter
type MockAIClient struct {
	SendChatRequestFunc func(req *domain.ChatRequest, chatID string) (*http.Response, error)
}

func (m *MockAIClient) SendChatRequest(req *domain.ChatRequest, chatID string) (*http.Response, error) {
	if m.SendChatRequestFunc != nil {
		return m.SendChatRequestFunc(req, chatID)
	}
	// Default mock behavior: return a dummy streaming response
	mockResponseContent := `id: chatcmpl-mockid
event: message
data: {"id":"chatcmpl-mockid","object":"chat.completion.chunk","created":1700000000,"model":"mock-model","choices":[{"index":0,"delta":{"content":"mocked "}}]}

id: chatcmpl-mockid
event: message
data: {"id":"chatcmpl-mockid","object":"chat.completion.chunk","created":1700000000,"model":"mock-model","choices":[{"index":0,"delta":{"content":"response"}}]}

id: chatcmpl-mockid
event: message
data: {"id":"chatcmpl-mockid","object":"chat.completion.chunk","created":1700000000,"model":"mock-model","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":"stop"}]}

id: chatcmpl-mockid
event: message
data: [DONE]

`
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "text/event-stream")
	recorder.WriteString(mockResponseContent)
	return recorder.Result(), nil
}

// MockTokenizer is a mock implementation of utils.Tokener
type MockTokenizer struct {
	InitFunc  func() error
	CountFunc func(text string) int
}

func (m *MockTokenizer) Init() error {
	if m.InitFunc != nil {
		return m.InitFunc()
	}
	return nil // Default successful init
}

func (m *MockTokenizer) Count(text string) int {
	if m.CountFunc != nil {
		return m.CountFunc(text)
	}
	return len(strings.Fields(text)) // Default behavior
}

// --- Helper Functions ---

func MockConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Port:    8080,
			Host:    "localhost",
			Debug:   true,
			Version: "test-v1.0",
		},
		Model: config.ModelConfig{
			Default: "mock-model",
		},
		Upstream: config.UpstreamConfig{
			Protocol: "https",
			Host:     "mock-z.ai",
			Token:    "mock-token",
		},
	}
}

// --- Tests ---

func TestNewRouter(t *testing.T) {
	cfg := MockConfig()
	mockAIClient := &MockAIClient{}
	mockTokenizer := &MockTokenizer{}
	router := api.NewRouter(cfg, mockAIClient, mockTokenizer)

	assert.NotNil(t, router)
	assert.IsType(t, &chi.Mux{}, router)
}

func TestRoutes(t *testing.T) {
	cfg := MockConfig()
	mockAIClient := &MockAIClient{}
	mockTokenizer := &MockTokenizer{}

	router := api.NewRouter(cfg, mockAIClient, mockTokenizer)

	// Test case for successful non-streaming chat completion
	t.Run("Chat Completions - Non-Streaming Success", func(t *testing.T) {
		mockAIClient.SendChatRequestFunc = func(req *domain.ChatRequest, chatID string) (*http.Response, error) {
			mockSSE := `data: {"data":{"delta_content":"Hello,","phase":"content"}}
data: {"data":{"delta_content":" world!","phase":"content"}}
data: {"data":{"done":true,"phase":"result"}}
data: [DONE]
`
			rr := httptest.NewRecorder()
			rr.Header().Set("Content-Type", "text/event-stream")
			rr.WriteString(mockSSE)
			return rr.Result(), nil
		}

		mockTokenizer.CountFunc = func(text string) int {
			if strings.Contains(text, "Hello, world!") {
				return 3
			}
			if strings.Contains(text, "hi") {
				return 1
			}
			return 0
		}

		body := `{"model": "test-model", "messages": [{"role": "user", "content": "hi"}], "stream": false}`
		req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader(body))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var respBody domain.ChatResponse
		err := json.Unmarshal(rr.Body.Bytes(), &respBody)
		require.NoError(t, err, "Failed to unmarshal response body")
		assert.Equal(t, "chat.completion", respBody.Object)
		require.Len(t, respBody.Choices, 1, "Should have exactly one choice")
		assert.Equal(t, "Hello, world!", respBody.Choices[0].Message.Content)
		assert.Equal(t, "assistant", respBody.Choices[0].Message.Role)
		assert.Equal(t, 1, respBody.Usage.PromptTokens)
		assert.Equal(t, 3, respBody.Usage.CompletionTokens)
	})

	// Test case for streaming chat completion
	t.Run("Chat Completions - Streaming Success", func(t *testing.T) {
		mockAIClient.SendChatRequestFunc = func(req *domain.ChatRequest, chatID string) (*http.Response, error) {
			mockSSE := `data: {"data":{"delta_content":"stream","phase":"content"}}
data: {"data":{"delta_content":"ing","phase":"content"}}
data: {"data":{"done":true}}
data: [DONE]
`
			rr := httptest.NewRecorder()
			rr.Header().Set("Content-Type", "text/event-stream")
			rr.WriteString(mockSSE)
			return rr.Result(), nil
		}

		mockTokenizer.CountFunc = func(text string) int {
			if strings.Contains(text, "hi") {
				return 1
			}
			if strings.Contains(text, "streaming") {
				return 2
			}
			return 0
		}

		body := `{"model": "test-model", "messages": [{"role": "user", "content": "hi"}], "stream": true}`
		req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader(body))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "text/event-stream", rr.Header().Get("Content-Type"))
		responseBody := rr.Body.String()
		assert.Contains(t, responseBody, "data: {")
		assert.Contains(t, responseBody, "\"content\":\"stream\"")
		assert.Contains(t, responseBody, "\"content\":\"ing\"")
		assert.Contains(t, responseBody, "data: [DONE]")
	})

	// Test case for AI client failure
	t.Run("Chat Completions - AI Client Error", func(t *testing.T) {
		mockAIClient.SendChatRequestFunc = func(req *domain.ChatRequest, chatID string) (*http.Response, error) {
			return nil, errors.New("internal AI error")
		}

		body := `{"model": "test-model", "messages": [{"role": "user", "content": "hi"}]}`
		req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader(body))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Failed to process request")
	})

	// Test case for invalid JSON
	t.Run("Chat Completions - Invalid JSON", func(t *testing.T) {
		// No AI client or tokenizer mock needed here as it fails earlier
		body := `{"messages": ` // Invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader(body))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid JSON request")
	})

	// Health and Models routes
	t.Run("Health and Models Routes", func(t *testing.T) {
		// These handlers don't depend on AIClienter or Tokener, so passing simple mocks is fine.
		router := api.NewRouter(cfg, &MockAIClient{}, &MockTokenizer{})

		// Health
		reqHealth := httptest.NewRequest(http.MethodGet, "/health", nil)
		rrHealth := httptest.NewRecorder()
		router.ServeHTTP(rrHealth, reqHealth)
		assert.Equal(t, http.StatusOK, rrHealth.Code)
		assert.Contains(t, rrHealth.Body.String(), "\"status\":\"ok\"")

		// Models
		reqModels := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
		rrModels := httptest.NewRecorder()
		router.ServeHTTP(rrModels, reqModels)
		assert.Equal(t, http.StatusOK, rrModels.Code)
		assert.Contains(t, rrModels.Body.String(), "\"object\":\"list\"")
	})
}