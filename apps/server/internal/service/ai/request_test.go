package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarazaex69/zik/apps/server/internal/config"
	"github.com/zarazaex69/zik/apps/server/internal/domain"
	"github.com/zarazaex69/zik/apps/server/internal/service/auth"
)

func TestFormatRequest(t *testing.T) {
	cfg := &config.Config{
		Model: config.ModelConfig{
			Default: "test-model",
		},
	}

	tests := []struct {
		name    string
		req     *domain.ChatRequest
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Simple text message",
			req: &domain.ChatRequest{
				Model: "gpt-4",
				Messages: []domain.Message{
					{Role: "user", Content: "Hello"},
				},
			},
			want: map[string]interface{}{
				"model": "gpt-4",
				"messages": []map[string]interface{}{
					{"role": "user", "content": "Hello"},
				},
				"stream": true,
			},
		},
		{
			name: "Default model",
			req: &domain.ChatRequest{
				Messages: []domain.Message{
					{Role: "user", Content: "Hello"},
				},
			},
			want: map[string]interface{}{
				"model": "test-model",
				"messages": []map[string]interface{}{
					{"role": "user", "content": "Hello"},
				},
				"stream": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatRequest(tt.req, cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check basic fields
			assert.Equal(t, tt.want["model"], got["model"])
			assert.Equal(t, tt.want["stream"], got["stream"])

			// Check messages
			gotMsgs, ok := got["messages"].([]map[string]interface{})
			assert.True(t, ok)
			wantMsgs, ok := tt.want["messages"].([]map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, len(wantMsgs), len(gotMsgs))

			for i, wantMsg := range wantMsgs {
				assert.Equal(t, wantMsg["role"], gotMsgs[i]["role"])
				assert.Equal(t, wantMsg["content"], gotMsgs[i]["content"])
			}
		})
	}
}

func TestUploadImage(t *testing.T) {
	// Mock auth server and upload server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/v1/auths/") {
			// Mock auth response
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    "user1",
				"token": "token1",
			})
			return
		}
		if strings.Contains(r.URL.Path, "/api/v1/files/") {
			// Mock upload response
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"id":       "file123",
				"filename": "image.png",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	// Parse test server URL
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

	// Reset auth service instance
	auth.NewService().ClearCache()

	t.Run("Upload valid base64 image", func(t *testing.T) {
		// 1x1 pixel transparent png base64
		base64Img := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="

		chatID := "chat123"
		got, err := UploadImage(base64Img, chatID, cfg)

		assert.NoError(t, err)
		assert.Equal(t, "file123_image.png", got)
	})

	t.Run("Skip non-base64", func(t *testing.T) {
		url := "http://example.com/image.png"
		got, err := UploadImage(url, "chat1", cfg)
		assert.NoError(t, err)
		assert.Equal(t, "", got)
	})
}
