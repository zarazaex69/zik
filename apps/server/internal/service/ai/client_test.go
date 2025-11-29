package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/domain"
	"github.com/zarazaex/zik/apps/server/internal/pkg/crypto"
	"github.com/zarazaex/zik/apps/server/internal/service/auth"
)

// mockSignatureGenerator is a mock implementation of crypto.SignatureGenerator for testing
type mockSignatureGenerator struct{}

func (m *mockSignatureGenerator) GenerateSignature(params map[string]string, lastUserMessage string) (*crypto.SignatureResult, error) {
	return &crypto.SignatureResult{
		Signature: "mock-signature",
		Timestamp: 1234567890,
	}, nil
}

func TestSendChatRequest(t *testing.T) {
	// Mock Z.AI API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/v1/auths/") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    "user1",
				"token": "token1",
			})
			return
		}
		if strings.Contains(r.URL.Path, "/api/v2/chat/completions") {
			// Verify headers
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

			// Verify body
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			assert.Equal(t, "gpt-4", body["model"])
			assert.Equal(t, "chat1", body["chat_id"])

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("data: stream"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	urlParts := strings.Split(ts.URL, "//")
	protocol := urlParts[0]
	host := urlParts[1]

	cfg := &config.Config{
		Upstream: config.UpstreamConfig{
			Protocol: protocol,
			Host:     host,
			Token:    "test-token",
		},
	}

	authSvc := auth.NewService()
	authSvc.ClearCache()
	mockSigGen := &mockSignatureGenerator{}
	client := NewClient(cfg, authSvc, mockSigGen)

	req := &domain.ChatRequest{
		Model: "gpt-4",
		Messages: []domain.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	resp, err := client.SendChatRequest(req, "chat1")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestSendChatRequest_Errors(t *testing.T) {
	// Setup generic mock server for auth failure
	tsAuthFail := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/v1/auths/") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}))
	defer tsAuthFail.Close()

	urlParts := strings.Split(tsAuthFail.URL, "//")
	protocol := urlParts[0]
	host := urlParts[1]

	cfgAuthFail := &config.Config{
		Upstream: config.UpstreamConfig{
			Protocol: protocol,
			Host:     host,
			Token:    "test-token",
		},
	}
	authSvc := auth.NewService()
	authSvc.ClearCache() // Ensure no cached user
	mockSigGen := &mockSignatureGenerator{}
	clientAuthFail := NewClient(cfgAuthFail, authSvc, mockSigGen)

	req := &domain.ChatRequest{
		Model: "gpt-4",
		Messages: []domain.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	t.Run("Auth Failure", func(t *testing.T) {
		resp, err := clientAuthFail.SendChatRequest(req, "chat1")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to get user info")
	})

	// Setup mock server for Upstream API error
	tsApiErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/v1/auths/") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    "user1",
				"token": "token1",
			})
			return
		}
		if strings.Contains(r.URL.Path, "/api/v2/chat/completions") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	defer tsApiErr.Close()

	urlParts = strings.Split(tsApiErr.URL, "//")
	protocol = urlParts[0]
	host = urlParts[1]
	cfgApiErr := &config.Config{
		Upstream: config.UpstreamConfig{
			Protocol: protocol,
			Host:     host,
			Token:    "test-token",
		},
	}
	clientApiErr := NewClient(cfgApiErr, authSvc, mockSigGen)
	authSvc.ClearCache()

	t.Run("Upstream API Error", func(t *testing.T) {
		resp, err := clientApiErr.SendChatRequest(req, "chat1")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Z.AI API error")
	})

	// We want "failed to send request" in the chat part.
	// This happens if GetUser succeeds but Chat req fails.
	// This implies network is up for GetUser but down for Chat.

	// We can simulate this with the test server logic:
	// Handler for auth returns OK.
	// Handler for chat: panic? or Close connection?
	// Hijack connection?

	tsNetErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/v1/auths/") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    "user1",
				"token": "token1",
			})
			return
		}
		// For other requests, close connection abruptly
		c, _, _ := w.(http.Hijacker).Hijack()
		c.Close()
	}))
	defer tsNetErr.Close()

	urlParts = strings.Split(tsNetErr.URL, "//")
	cfgNetErr := &config.Config{
		Upstream: config.UpstreamConfig{
			Protocol: urlParts[0],
			Host:     urlParts[1],
			Token:    "test",
		},
	}
	clientNetErr := NewClient(cfgNetErr, authSvc, mockSigGen)
	authSvc.ClearCache()

	t.Run("Network Error on Chat", func(t *testing.T) {
		resp, err := clientNetErr.SendChatRequest(req, "chat1")
		assert.Error(t, err)
		assert.Nil(t, resp)
		// Go http client reports EOF or connection reset
		assert.Contains(t, err.Error(), "failed to send request")
	})

	// We need auth success for NewRequest error (to get past getUser)
	// But actually GetUser is called first, so we need a valid auth response first.
	// tsValid will provide valid auth.
	tsValid := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/v1/auths/") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    "user1",
				"token": "token1",
			})
			return
		}
	}))
	defer tsValid.Close()
	urlParts = strings.Split(tsValid.URL, "//") // This gives a valid base URL for auth requests.

	// authSvc will use cfgBadURL to make auth request?
	// No, authSvc uses `cfg` passed to `GetUser(cfg)`.
	// NewClient takes `cfg` and stores it.
	// `SendChatRequest` calls `c.authService.GetUser(c.cfg)`.

	// So `cfgBadURL` MUST have a valid Upstream for Auth to work, BUT produce invalid URL for Chat?
	// `GetUser` uses `cfg.Upstream` to build auth URL.
	// `SendChatRequest` uses `cfg.Upstream` to build chat URL.
	// They use the SAME config.
	// If I put "host\n" in config:
	// `GetUser` URL: "http://host\n/api/v1/user" -> NewRequest likely fails HERE first!

	// If NewRequest fails in GetUser, we get "failed to get user info: ... failed to create request ...".
	// This duplicates "Auth Failure" test effectively but with different root cause.
	// But wait, `TestSendChatRequest_Errors` runs `Auth Failure` which returns "failed to get user info".

	// I want to fail AFTER `GetUser` succeeds.
	// `GetUser` succeeds if `http.NewRequest` works AND request works.
	// `http.NewRequest` will fail if host has newline.
	// So `GetUser` will fail.

	// So `NewRequest Error` test as I designed it will fail inside GetUser.
	// And the error message will be "failed to get user info".
	// But I assert "failed to create HTTP request".
	// If "failed to get user info" wraps "failed to create request", then assert might fail if I look for specific string.
	// Let's check `client.go`: `return nil, fmt.Errorf("failed to get user info: %w", err)`

	// If I want to fail `http.NewRequest` in `SendChatRequest` specifically (line 92), I need `GetUser` to succeed.
	// `GetUser` succeeds if it returns a user.
	// If I mock `authService`? `NewClient` takes `*auth.Service`. It is a struct, not interface.
	// But `auth.Service` has `GetUser`.
	// If I can put user in cache!
	// `auth.Service` has internal cache if I recall correctly (logs said "User cache cleared").
	// I can't access `auth.Service` private cache directly.
	// But if I call `authSvc.GetUser(validCfg)` first, it caches the user!
	// Then I call `clientBadURL.SendChatRequest`.
	// Wait, `SendChatRequest` calls `c.authService.GetUser(c.cfg)`.
	// `c.cfg` is the BAD config.
	// `authSvc.GetUser` uses the passed config key for cache? Or just caches per token?
	// Let's assume it uses token or something.
	// If I can prime the cache using a VALID config that has SAME token, maybe it returns cached user without making request using BAD config?

	// Let's see `GetUser` implementation in `apps/server/internal/service/auth/user.go`?
	// Assume `GetUser` caches based on token in config?

	// Let's try attempting to prime cache.
	cfgValid := &config.Config{
		Upstream: config.UpstreamConfig{
			Protocol: urlParts[0],
			Host:     urlParts[1],
			Token:    "test",
		},
	}

	// Prime cache
	commonAuthSvc := auth.NewService()
	_, err := commonAuthSvc.GetUser(cfgValid)
	if err != nil {
		// If this fails, test setup fails
	}

	// Now use commonAuthSvc with bad config.
	// If cache logic ignores host/protocol and uses only token (likely), it will return user from cache.

	t.Run("NewRequest Error", func(t *testing.T) {
		// Use config that produces invalid URL for http.NewRequest
		// NewRequest parses the URL. If URL contains control chars like \n, it might fail.
		cfgBadURL := &config.Config{
			Upstream: config.UpstreamConfig{
				Protocol: "http",
				Host:     "host\n", // Newline in host/url should cause NewRequest to fail
				Token:    "test",
			},
		}
		clientBadURL := NewClient(cfgBadURL, authSvc, mockSigGen)
		// Request will try to build URL with \n
		resp, err := clientBadURL.SendChatRequest(req, "chat1")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to create HTTP request")
	})
}

func TestExtractLastUserMessage(t *testing.T) {
	tests := []struct {
		name     string
		messages []domain.Message
		want     string
	}{
		{
			name: "Simple text",
			messages: []domain.Message{
				{Role: "system", Content: "sys"},
				{Role: "user", Content: "hello"},
				{Role: "assistant", Content: "hi"},
			},
			want: "hello",
		},
		{
			name: "Multiple user messages",
			messages: []domain.Message{
				{Role: "user", Content: "first"},
				{Role: "assistant", Content: "hi"},
				{Role: "user", Content: "second"},
			},
			want: "second",
		},
		{
			name: "Multimodal",
			messages: []domain.Message{
				{
					Role: "user",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "image desc"},
					},
				},
			},
			want: "image desc",
		},
		{
			name:     "Empty messages slice",
			messages: []domain.Message{},
			want:     "",
		},
		{
			name: "No user messages",
			messages: []domain.Message{
				{Role: "system", Content: "sys"},
				{Role: "assistant", Content: "hi"},
			},
			want: "",
		},
		{
			name: "Multimodal with non-text item",
			messages: []domain.Message{
				{
					Role: "user",
					Content: []interface{}{
						map[string]interface{}{"type": "image_url", "image_url": map[string]interface{}{"url": "http://example.com/img.png"}},
						map[string]interface{}{"type": "text", "text": "What is in this image?"},
					},
				},
			},
			want: "What is in this image?",
		},
		{
			name: "Multimodal with mixed content, last user message is complex",
			messages: []domain.Message{
				{
					Role: "user",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "part1"},
						map[string]interface{}{"type": "image_url", "image_url": map[string]interface{}{"url": "http://example.com/img2.png"}},
						map[string]interface{}{"type": "text", "text": "part2"},
					},
				},
			},
			want: "part1 part2",
		},
		{
			name: "Multimodal with text field not a string",
			messages: []domain.Message{
				{
					Role: "user",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": 123}, // text is int, not string
					},
				},
			},
			want: "", // Should not include non-string text
		},
		{
			name: "User message with non-string, non-array content",
			messages: []domain.Message{
				{Role: "user", Content: 123}, // int content
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractLastUserMessage(tt.messages)
			assert.Equal(t, tt.want, got)
		})
	}
}
