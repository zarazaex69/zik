package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zarazaex/zik/apps/server/internal/api/middleware"
	"github.com/zarazaex/zik/apps/server/internal/pkg/logger"
)

type mockFlusherResponseWriter struct {
	*httptest.ResponseRecorder
	flushed bool
}

func (m *mockFlusherResponseWriter) Flush() {
	m.flushed = true
}

func NewMockFlusherResponseWriter() *mockFlusherResponseWriter {
	return &mockFlusherResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
	}
}

func TestLoggerMiddleware(t *testing.T) {
	// 1. Setup: Redirect logger output to a buffer
	var logBuffer bytes.Buffer
	// Use debug=false to get JSON output which is easier to parse
	logger.InitWithWriter(false, &logBuffer)

	// 2. Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("test response"))
	})

	// 3. Wrap handler with Logger middleware
	loggedHandler := middleware.Logger(testHandler)

	// 4. Create a request and response recorder
	req := httptest.NewRequest(http.MethodPost, "/test-path", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rr := httptest.NewRecorder()

	// 5. Serve the request
	loggedHandler.ServeHTTP(rr, req)

	// 6. Assert response is correct
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, "test response", rr.Body.String())

	// 7. Assert log output is correct
	logOutput := logBuffer.String()
	t.Log(logOutput) // a good practice to see the output when debugging tests
	require.True(t, strings.Contains(logOutput, "HTTP request"), "Log should contain the message")
	assert.True(t, strings.Contains(logOutput, `"method":"POST"`), "Log should contain the correct method")
	assert.True(t, strings.Contains(logOutput, `"path":"/test-path"`), "Log should contain the correct path")
	assert.True(t, strings.Contains(logOutput, `"status":202`), "Log should contain the correct status code")
	assert.True(t, strings.Contains(logOutput, `"remote_addr":"127.0.0.1:12345"`), "Log should contain the correct remote address")
}

func TestResponseWriter_Flush(t *testing.T) {
	t.Run("with flusher support", func(t *testing.T) {
		// 1. Create a mock response writer that implements Flusher
		mockWriter := NewMockFlusherResponseWriter()

		// 2. Wrap it in our ResponseWriter
		wrapped := middleware.NewResponseWriter(mockWriter)

		// 3. Call Flush
		wrapped.Flush()

		// 4. Assert that the underlying Flush was called
		assert.True(t, mockWriter.flushed, "underlying Flush should have been called")
	})

	t.Run("without flusher support", func(t *testing.T) {
		// 1. Create a standard response recorder (does not implement Flusher)
		recorder := httptest.NewRecorder()

		// 2. Wrap it
		wrapped := middleware.NewResponseWriter(recorder)

		        // 3. Call Flush (should not panic)
				// The wrapped writer itself is what implements Flusher.
				// The check for the underlying writer supporting it is inside the Flush method.
				assert.NotPanics(t, func() {
					wrapped.Flush()
				})	})
}