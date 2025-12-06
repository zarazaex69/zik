package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New(10 * time.Second)

	if client == nil {
		t.Fatal("New() returned nil client")
	}

	if client.Timeout != 10*time.Second {
		t.Errorf("New() timeout = %v, want %v", client.Timeout, 10*time.Second)
	}

	if client.Transport == nil {
		t.Error("New() client.Transport is nil")
	}
}

func TestNew_DifferentTimeouts(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{"1 second", 1 * time.Second},
		{"5 seconds", 5 * time.Second},
		{"30 seconds", 30 * time.Second},
		{"1 minute", 1 * time.Minute},
		{"zero timeout", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(tt.timeout)

			if client.Timeout != tt.timeout {
				t.Errorf("New(%v) timeout = %v, want %v", tt.timeout, client.Timeout, tt.timeout)
			}
		})
	}
}

func TestNew_TransportConfiguration(t *testing.T) {
	client := New(10 * time.Second)

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("client.Transport is not *http.Transport")
	}

	// Verify transport settings
	if !transport.ForceAttemptHTTP2 {
		t.Error("ForceAttemptHTTP2 should be true")
	}

	if transport.MaxIdleConns != 100 {
		t.Errorf("MaxIdleConns = %d, want 100", transport.MaxIdleConns)
	}

	if transport.IdleConnTimeout != 90*time.Second {
		t.Errorf("IdleConnTimeout = %v, want 90s", transport.IdleConnTimeout)
	}

	if transport.TLSHandshakeTimeout != 10*time.Second {
		t.Errorf("TLSHandshakeTimeout = %v, want 10s", transport.TLSHandshakeTimeout)
	}

	if transport.ExpectContinueTimeout != 1*time.Second {
		t.Errorf("ExpectContinueTimeout = %v, want 1s", transport.ExpectContinueTimeout)
	}
}

func TestNew_HTTPRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	client := New(10 * time.Second)

	// Make request
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("client.Get() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("response status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestNew_ContextTimeout(t *testing.T) {
	// Create slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(100 * time.Millisecond)

	// Request should timeout
	_, err := client.Get(server.URL)
	if err == nil {
		t.Error("client.Get() should timeout but didn't")
	}
}

func TestNew_WithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(10 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext() error = %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("client.Do() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("response status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestNew_CancelledContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(10 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext() error = %v", err)
	}

	_, err = client.Do(req)
	if err == nil {
		t.Error("client.Do() should fail with cancelled context")
	}
}

func TestNew_MultipleClients(t *testing.T) {
	client1 := New(5 * time.Second)
	client2 := New(10 * time.Second)

	if client1.Timeout == client2.Timeout {
		t.Error("Different clients should have different timeouts")
	}

	// Verify they are independent instances
	if client1 == client2 {
		t.Error("New() should return different client instances")
	}
}

func TestNew_ProxyFromEnvironment(t *testing.T) {
	client := New(10 * time.Second)

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("client.Transport is not *http.Transport")
	}

	// Verify proxy function is set
	if transport.Proxy == nil {
		t.Error("Proxy should be set to http.ProxyFromEnvironment")
	}
}

func TestNew_CustomResolver(t *testing.T) {
	client := New(10 * time.Second)

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("client.Transport is not *http.Transport")
	}

	// Verify DialContext is set (indicates custom resolver)
	if transport.DialContext == nil {
		t.Error("DialContext should be set for custom DNS resolver")
	}
}

func TestNew_HTTP2Support(t *testing.T) {
	client := New(10 * time.Second)

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("client.Transport is not *http.Transport")
	}

	if !transport.ForceAttemptHTTP2 {
		t.Error("HTTP/2 support should be enabled")
	}
}

func TestNew_ConnectionPooling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(10 * time.Second)

	// Make multiple requests to test connection pooling
	for i := 0; i < 5; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
		resp.Body.Close()
	}
}
