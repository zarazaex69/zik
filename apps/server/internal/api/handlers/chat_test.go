package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zarazaex/zik/apps/server/internal/api/handlers"
	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/domain"
)

// --- Mock Implementations ---

// MockAuthService is a mock implementation of auth.AuthServicer for testing
type MockAuthService struct {
	GetUserFunc func(cfg *config.Config) (*domain.User, error)
}

func (m *MockAuthService) GetUser(cfg *config.Config) (*domain.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(cfg)
	}
	return &domain.User{ID: "mock-user-id", Token: "mock-user-token"}, nil
}

func (m *MockAuthService) ClearCache() {
	// No-op for mock
}

// MockAIClient is a mock implementation of ai.AIClienter for testing
type MockAIClient struct {
	SendChatRequestFunc func(req *domain.ChatRequest, chatID string) (*http.Response, error)
}

func (m *MockAIClient) SendChatRequest(req *domain.ChatRequest, chatID string) (*http.Response, error) {
	if m.SendChatRequestFunc != nil {
		return m.SendChatRequestFunc(req, chatID)
	}
	// Default mock behavior for streaming
	mockResponseContent := `data: {"data":{"delta_content":"mocked "}}` + "\n\n" +
		`data: {"data":{"delta_content":"response."}}` + "\n\n" +
		`data: {"data":{"done":true}}` + "\n\n" +
		`data: [DONE]` + "\n"
	rr := httptest.NewRecorder()
	rr.Header().Set("Content-Type", "text/event-stream")
	rr.WriteString(mockResponseContent)
	return rr.Result(), nil
}

// MockTokenizer is a mock implementation of utils.Tokenizer
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
	// Default behavior: return 1 token per word, or 0 if empty
	if text == "" {
		return 0
	}
	return len(strings.Fields(text))
}

// --- Helper Functions ---

func MockConfig() *config.Config {
	return &config.Config{
		Upstream: config.UpstreamConfig{
			Protocol: "https",
			Host:     "mock-z.ai",
			Token:    "test-token",
		},
		Model: config.ModelConfig{
			Default: "gpt-4",
		},
	}
}

// --- Tests ---

func TestChatCompletions(t *testing.T) {
	cfg := MockConfig()
	mockAIClient := &MockAIClient{}
	mockTokenizer := &MockTokenizer{}

	// Instantiate the real ai.Client with the mock auth service
	// This is needed because handlers.ChatCompletions takes ai.AIClienter,
	// but ai.NewClient takes auth.AuthServicer.
	// So, we need to adapt our mockAIClient to be compatible with ai.AIClienter.
	// The aiClient passed to ChatCompletions is actually a concrete *MockAIClient.
	// So it works by polymorphism.

	// The problem is that the original TestChatCompletions used httptest.NewServer
	// to mock the *external* dependencies (auth and AI).
	// Now with interfaces, we mock the *internal* dependencies directly.

	handler := handlers.ChatCompletions(cfg, mockAIClient, mockTokenizer)

	t.Run("Non-streaming request", func(t *testing.T) {
		// Override AI client mock to return a non-streaming response
		mockAIClient.SendChatRequestFunc = func(req *domain.ChatRequest, chatID string) (*http.Response, error) {
			// Simulate a successful, non-streaming SSE response from Z.AI
			mockSSE := `data: {"data":{"delta_content":"Hello,","phase":"content"}}` + "\n\n" +
				`data: {"data":{"delta_content":" world!","phase":"content"}}` + "\n\n" +
				`data: {"data":{"done":true,"phase":"result"}}` + "\n\n" +
				`data: [DONE]` + "\n"
			rr := httptest.NewRecorder()
			rr.Header().Set("Content-Type", "text/event-stream")
			rr.WriteString(mockSSE)
			return rr.Result(), nil
		}

		mockTokenizer.CountFunc = func(text string) int {
			if text == "Hello, world!" {
				return 3 // "Hello,", " world!", "!"
			}
			if strings.Contains(text, "Hi") {
				return 1 // "Hi"
			}
			return 0
		}

		reqBody := domain.ChatRequest{
			Model: "gpt-4",
			Messages: []domain.Message{
				{Role: "user", Content: "Hi"},
			},
			Stream: false,
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var chatResp domain.ChatResponse
		err := json.NewDecoder(resp.Body).Decode(&chatResp)
		require.NoError(t, err)
		assert.NotEmpty(t, chatResp.ID)
		assert.Equal(t, "chat.completion", chatResp.Object)
		require.Len(t, chatResp.Choices, 1)
		assert.Equal(t, "assistant", chatResp.Choices[0].Message.Role)
		assert.Equal(t, "Hello, world!", chatResp.Choices[0].Message.Content)
		assert.Equal(t, "stop", *chatResp.Choices[0].FinishReason)
		assert.NotNil(t, chatResp.Usage)
		assert.Equal(t, 1, chatResp.Usage.PromptTokens)
		assert.Equal(t, 3, chatResp.Usage.CompletionTokens)
		assert.Equal(t, 4, chatResp.Usage.TotalTokens)
	})

	t.Run("Streaming request", func(t *testing.T) {
		// Reset to default streaming mock behavior
		mockAIClient.SendChatRequestFunc = nil

		mockTokenizer.CountFunc = func(text string) int {
			if strings.Contains(text, "Hi") {
				return 1 // "Hi"
			}
			if strings.Contains(text, "streaming") {
				return 2 // "streaming"
			}
			return 0
		}

		reqBody := domain.ChatRequest{
			Model: "gpt-4",
			Messages: []domain.Message{
				{Role: "user", Content: "Hi"},
			},
			Stream: true,
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))

		bodyStr := w.Body.String()
		assert.Contains(t, bodyStr, "data: {")
		assert.Contains(t, bodyStr, "\"content\":\"mocked \"")
		assert.Contains(t, bodyStr, "\"content\":\"response.\"")
		assert.Contains(t, bodyStr, "data: [DONE]")
		// Usage is not included by default in streaming (only if StreamOpts.IncludeUsage is true)
	})

	t.Run("Streaming request with usage include", func(t *testing.T) {
		mockAIClient.SendChatRequestFunc = func(req *domain.ChatRequest, chatID string) (*http.Response, error) {
			mockSSE := `data: {"data":{"delta_content":"stream","phase":"content"}}` + "\n\n" +
				`data: {"data":{"delta_content":"ing","phase":"content"}}` + "\n\n" +
				`data: {"data":{"done":true}}` + "\n\n" +
				`data: [DONE]` + "\n"
			rr := httptest.NewRecorder()
			rr.Header().Set("Content-Type", "text/event-event") // Correct SSE
			rr.WriteString(mockSSE)
			return rr.Result(), nil
		}

		mockTokenizer.CountFunc = func(text string) int {
			if strings.Contains(text, "Hi") {
				return 1
			}
			if strings.Contains(text, "streaming") {
				return 2
			}
			return 0
		}

		reqBody := domain.ChatRequest{
			Model: "gpt-4",
			Messages: []domain.Message{
				{Role: "user", Content: "Hi"},
			},
			Stream:    true,
			StreamOpts: &domain.StreamOptions{IncludeUsage: true},
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))
		bodyStr := w.Body.String()

		assert.Contains(t, bodyStr, "\"usage\":{\"prompt_tokens\":1,\"completion_tokens\":2,\"total_tokens\":3}")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader([]byte("invalid")))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("Validation failure", func(t *testing.T) {
		reqBody := domain.ChatRequest{
			Model: "gpt-4",
			// Missing messages field, which is required
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "validation failed")
	})

	t.Run("AI client returns error", func(t *testing.T) {
		mockAIClient.SendChatRequestFunc = func(req *domain.ChatRequest, chatID string) (*http.Response, error) {
			return nil, errors.New("simulated AI client error")
		}

		reqBody := domain.ChatRequest{
			Model: "gpt-4",
			Messages: []domain.Message{
				{Role: "user", Content: "Hi"},
			},
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Failed to process request")
	})

	t.Run("Streaming not supported (non-flusher writer)", func(t *testing.T) {
		// Mock a response writer that does NOT implement http.Flusher
		originalRW := httptest.NewRecorder() // This is what httptest.NewRecorder returns, which is not a flusher by default.

		reqBody := domain.ChatRequest{
			Model: "gpt-4",
			Messages: []domain.Message{
				{Role: "user", Content: "Hi"},
			},
			Stream: true, // Request streaming
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(bodyBytes))

		// Override AI client to return error to simulate flusher check failure
		mockAIClient.SendChatRequestFunc = func(req *domain.ChatRequest, chatID string) (*http.Response, error) {
			return nil, errors.New("simulated AI client error")
		}

		handler(originalRW, req)

		assert.Equal(t, http.StatusInternalServerError, originalRW.Result().StatusCode)
		assert.Contains(t, originalRW.Body.String(), "Failed to process request")
	})
}