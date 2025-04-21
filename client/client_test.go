package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestRetryLogic(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)

		if count <= 2 {
			http.Error(w, "temporary failure", http.StatusInternalServerError)
			return
		}
		_, err := fmt.Fprintln(w, `{"message": "success"}`)
		if err != nil {
			return
		}
	}))
	defer server.Close()

	c := New(
		server.URL,
		rate.Every(10*time.Millisecond),
		1,
		5*time.Second,
		3,
		10*time.Millisecond,
		map[string]string{"Content-Type": "application/json"},
	)

	resp, err := c.DoRequest(context.Background(), "GET", "/", nil)
	if err != nil {
		t.Fatalf("expected successful request, got error: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}

	if atomic.LoadInt32(&callCount) != 3 {
		t.Errorf("expected 3 attempts, got %d", callCount)
	}
}
