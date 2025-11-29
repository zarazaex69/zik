package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zarazaex/zik/apps/server/internal/config"
)

func TestGetUser(t *testing.T) {
	// Mock Auth API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer invalid" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    "user123",
			"name":  "Test User",
			"token": "token123",
		})
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

	s := NewService()
	s.ClearCache()

	t.Run("Get user from API", func(t *testing.T) {
		user, err := s.GetUser(cfg)
		assert.NoError(t, err)
		assert.Equal(t, "user123", user.ID)
		assert.Equal(t, "test-token", user.Token)
	})

	t.Run("Get user from cache", func(t *testing.T) {
		// First call to populate cache
		_, err := s.GetUser(cfg)
		assert.NoError(t, err)

		// Modify config to break API call (to prove cache is used)
		badCfg := &config.Config{
			Upstream: config.UpstreamConfig{
				Protocol: "http:",
				Host:     "invalid-host",
				Token:    "test-token",
			},
		}

		user, err := s.GetUser(badCfg)
		assert.NoError(t, err)
		assert.Equal(t, "user123", user.ID)
	})

	t.Run("API error", func(t *testing.T) {
		s.ClearCache()
		badCfg := &config.Config{
			Upstream: config.UpstreamConfig{
				Protocol: protocol,
				Host:     host,
				Token:    "invalid",
			},
		}

		_, err := s.GetUser(badCfg)
		assert.Error(t, err)
	})
}

func TestClearCache(t *testing.T) {
	s := NewService()
	s.cache["token"] = &cachedUser{cachedAt: time.Now()}

	s.ClearCache()

	assert.Empty(t, s.cache)
}
