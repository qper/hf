package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBuildServerRegistersHealthRoute(t *testing.T) {
	srv := newServer()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestServerStarts(t *testing.T) {
	go main()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/healthz")
	if err != nil {
		t.Fatalf("expected health endpoint to respond: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			t.Fatalf("failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
