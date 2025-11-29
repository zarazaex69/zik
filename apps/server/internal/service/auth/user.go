package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/zarazaex/zik/apps/server/internal/config"
	"github.com/zarazaex/zik/apps/server/internal/domain"
	"github.com/zarazaex/zik/apps/server/internal/pkg/logger"
)

// Service handles user authentication with caching
type Service struct {
	cache map[string]*cachedUser
	mutex sync.RWMutex
}

type cachedUser struct {
	user     *domain.User
	cachedAt time.Time
}

var (
	instance *Service
	once     sync.Once
)

// NewService creates a new auth service instance
func NewService() *Service {
	once.Do(func() {
		instance = &Service{
			cache: make(map[string]*cachedUser),
		}
	})
	return instance
}

// GetUser retrieves user information from Z.AI API with caching
func (s *Service) GetUser(cfg *config.Config) (*domain.User, error) {
	token := cfg.Upstream.Token

	// Check cache for authenticated users
	if token != "" {
		s.mutex.RLock()
		cached, exists := s.cache[token]
		s.mutex.RUnlock()

		// Return cached user if still valid (30 min TTL)
		if exists && time.Since(cached.cachedAt) < 30*time.Minute {
			logger.Debug().Str("user_id", cached.user.ID).Msg("Using cached user info")
			return cached.user, nil
		}
	}

	// Fetch from upstream API
	url := fmt.Sprintf("%s//%s/api/v1/auths/", cfg.Upstream.Protocol, cfg.Upstream.Host)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth request: %w", err)
	}

	// Add headers
	for k, v := range cfg.GetUpstreamHeaders() {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add authorization for authenticated users
	if !cfg.Upstream.Anonymous && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth API returned status %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode auth response: %w", err)
	}

	// Extract user info
	userID := getStringField(result, "id")
	userName := getStringField(result, "name")
	userToken := getStringField(result, "token")

	// Use provided token for authenticated users
	if !cfg.Upstream.Anonymous {
		userToken = token
	}

	user := &domain.User{
		ID:    userID,
		Token: userToken,
	}

	// Cache authenticated users
	if userToken != "" && userID != "" {
		s.mutex.Lock()
		s.cache[userToken] = &cachedUser{
			user:     user,
			cachedAt: time.Now(),
		}
		s.mutex.Unlock()
		logger.Info().Str("user_id", userID).Str("name", userName).Msg("User authenticated")
	}

	return user, nil
}

// ClearCache clears the user cache
func (s *Service) ClearCache() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cache = make(map[string]*cachedUser)
	logger.Info().Msg("User cache cleared")
}

func getStringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
