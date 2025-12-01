package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	tests := []struct {
		name           string
		requests       int
		expectedStatus int
		description    string
	}{
		{
			name:           "within_limit",
			requests:       3,
			expectedStatus: http.StatusOK,
			description:    "Requests within limit should succeed",
		},
		{
			name:           "at_limit",
			requests:       5,
			expectedStatus: http.StatusOK,
			description:    "Request at exact limit should succeed",
		},
		{
			name:           "exceed_limit",
			requests:       6,
			expectedStatus: http.StatusTooManyRequests,
			description:    "Request exceeding limit should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new limiter for each test to avoid interference
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			limited := RateLimit(5)(handler)

			var lastStatus int
			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "192.168.1.1:1234" // Same IP for all requests
				w := httptest.NewRecorder()

				limited.ServeHTTP(w, req)
				lastStatus = w.Code
			}

			if lastStatus != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.description, tt.expectedStatus, lastStatus)
			}
		})
	}
}

func TestRateLimitDifferentIPs(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limited := RateLimit(2)(handler)

	// First IP makes 2 requests (at limit)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		limited.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d from IP1 should succeed, got status %d", i+1, w.Code)
		}
	}

	// Second IP should still be able to make requests
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.2:1234"
	w := httptest.NewRecorder()
	limited.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Request from different IP should succeed, got status %d", w.Code)
	}
}

func TestRateLimitRefill(t *testing.T) {
	// This test would require waiting for time to pass
	// Skipping in unit tests, but demonstrates the concept
	t.Skip("Skipping time-dependent test")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limited := RateLimit(2)(handler)

	// Exhaust limit
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		limited.ServeHTTP(w, req)
	}

	// Next request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()
	limited.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected rate limit, got status %d", w.Code)
	}

	// Wait for refill
	time.Sleep(61 * time.Second)

	// Should work again
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w = httptest.NewRecorder()
	limited.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("After refill, request should succeed, got status %d", w.Code)
	}
}
