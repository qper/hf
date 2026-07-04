package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os/exec"
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

func TestStopProcessByPID(t *testing.T) {
	cmd := exec.Command("sh", "-c", "sleep 30")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child process: %v", err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	if err := stopProcessByPID(cmd.Process.Pid); err != nil {
		t.Fatalf("stop child process: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatalf("expected child process to exit with an error after signal")
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("child process did not stop in time")
	}
}

func TestServerStarts(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen for test server: %v", err)
	}
	addr := listener.Addr().String()
	_ = listener.Close()

	srv := newServer()
	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- srv.Start(addr)
	}()

	deadline := time.Now().Add(5 * time.Second)
	var resp *http.Response
	var reqErr error
	for time.Now().Before(deadline) {
		resp, reqErr = http.Get(fmt.Sprintf("http://%s/healthz", addr))
		if reqErr == nil && resp.StatusCode == http.StatusOK {
			break
		}
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}

	if reqErr != nil {
		t.Fatalf("expected health endpoint to respond: %v", reqErr)
	}
	if resp == nil {
		t.Fatalf("expected health endpoint response, got nil")
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			t.Fatalf("failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		t.Fatalf("shutdown server: %v", err)
	}

	select {
	case err := <-serverErrCh:
		if err != nil && err != http.ErrServerClosed {
			t.Fatalf("server returned unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("server did not stop promptly")
	}
}
