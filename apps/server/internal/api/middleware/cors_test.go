package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap with CORS middleware
	corsHandler := CORS(handler)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   map[string]string
	}{
		{
			name:           "GET request with CORS headers",
			method:         "GET",
			expectedStatus: http.StatusOK,
			checkHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization, X-Requested-With",
			},
		},
		{
			name:           "POST request with CORS headers",
			method:         "POST",
			expectedStatus: http.StatusOK,
			checkHeaders: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		},
		{
			name:           "OPTIONS preflight request",
			method:         "OPTIONS",
			expectedStatus: http.StatusOK,
			checkHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Max-Age":       "86400",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			corsHandler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("CORS() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			for header, expectedValue := range tt.checkHeaders {
				actualValue := w.Header().Get(header)
				if actualValue != expectedValue {
					t.Errorf("CORS() header %s = %s, want %s", header, actualValue, expectedValue)
				}
			}
		})
	}
}

func TestCORS_PreflightDoesNotCallHandler(t *testing.T) {
	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	corsHandler := CORS(handler)

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()

	corsHandler.ServeHTTP(w, req)

	if handlerCalled {
		t.Error("CORS() OPTIONS request should not call the wrapped handler")
	}

	if w.Code != http.StatusOK {
		t.Errorf("CORS() OPTIONS status = %d, want %d", w.Code, http.StatusOK)
	}
}
