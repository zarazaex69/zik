package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// RateLimiter implements a simple token bucket rate limiter per IP
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // requests per minute
	burst    int           // max burst size
	cleanup  time.Duration // cleanup interval for old visitors
}

// visitor represents a single IP's rate limit state
type visitor struct {
	limiter  *tokenBucket
	lastSeen time.Time
}

// tokenBucket implements a simple token bucket algorithm
type tokenBucket struct {
	tokens    int
	maxTokens int
	refillAt  time.Time
	mu        sync.Mutex
}

// NewRateLimiter creates a new rate limiter with specified requests per minute
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     requestsPerMinute,
		burst:    requestsPerMinute, // Allow burst up to rate limit
		cleanup:  5 * time.Minute,   // Clean up visitors not seen for 5 minutes
	}

	// Start cleanup goroutine to prevent memory leaks from inactive IPs
	go rl.cleanupVisitors()

	return rl
}

// allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	v, exists := rl.visitors[ip]
	if !exists {
		// Create new visitor with full token bucket
		v = &visitor{
			limiter: &tokenBucket{
				tokens:    rl.burst,
				maxTokens: rl.burst,
				refillAt:  time.Now().Add(time.Minute),
			},
			lastSeen: time.Now(),
		}
		rl.visitors[ip] = v
	}
	rl.mu.Unlock()

	v.lastSeen = time.Now()
	return v.limiter.take()
}

// take attempts to take a token from the bucket
// Returns true if token was available, false if rate limited
func (tb *tokenBucket) take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()

	// Refill bucket if minute has passed
	if now.After(tb.refillAt) {
		tb.tokens = tb.maxTokens
		tb.refillAt = now.Add(time.Minute)
	}

	// Check if tokens available
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// cleanupVisitors removes visitors that haven't been seen recently to prevent memory leaks
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.Sub(v.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit returns a middleware that limits requests per IP
func RateLimit(requestsPerMinute int) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(requestsPerMinute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract IP from request, considering X-Forwarded-For for proxies
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.Header.Get("X-Real-IP")
			}
			if ip == "" {
				ip = r.RemoteAddr
			}

			// Check rate limit
			if !limiter.allow(ip) {
				log.Warn().
					Str("ip", ip).
					Str("path", r.URL.Path).
					Msg("Rate limit exceeded")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":{"message":"Rate limit exceeded. Maximum 30 requests per minute.","type":"rate_limit_error","code":429}}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
