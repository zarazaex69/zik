package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecovery(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedStatus int
		shouldPanic    bool
	}{
		{
			name: "normal request without panic",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			},
			expectedStatus: http.StatusOK,
			shouldPanic:    false,
		},
		{
			name: "request with panic",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
			shouldPanic:    true,
		},
		{
			name: "request with nil panic",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic(nil)
			},
			expectedStatus: http.StatusInternalServerError,
			shouldPanic:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recoveryHandler := Recovery(tt.handler)

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Should not panic even if handler panics
			defer func() {
				if r := recover(); r != nil {
					t.Error("Recovery() middleware did not catch panic")
				}
			}()

			recoveryHandler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Recovery() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.shouldPanic {
				// Check that error response is JSON
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Recovery() Content-Type = %s, want application/json", contentType)
				}

				// Check that body contains error
				body := w.Body.String()
				if body == "" {
					t.Error("Recovery() returned empty body for panic")
				}
			}
		})
	}
}

func TestRecovery_DoesNotAffectNormalFlow(t *testing.T) {
	expectedBody := "test response"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(expectedBody))
	})

	recoveryHandler := Recovery(handler)

	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	recoveryHandler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Recovery() status = %d, want %d", w.Code, http.StatusCreated)
	}

	if w.Body.String() != expectedBody {
		t.Errorf("Recovery() body = %s, want %s", w.Body.String(), expectedBody)
	}
}
