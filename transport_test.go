package ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestTransport(t *testing.T) {
	s := NewStore()
	s.Register("test", 10, 10)
	tr := NewTransport(s)

	ctx := context.Background()
	ctx = ToContext(ctx, "test")
	limiter, ok := tr.fromContext(ctx)
	if !ok {
		t.Fatalf("expected true but got false")
	}
	if limiter == nil {
		t.Fatalf("expected limiter instance but got nil")
	}
}

func TestRoundTrip(t *testing.T) {
	testRateLimiterKey := "test"
	rateLimitPerSecond := 2
	s := NewStore()
	s.Register(testRateLimiterKey, rateLimitPerSecond, 1)
	tr := NewTransport(s)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, test")
	}))
	defer ts.Close()

	httpCli := &http.Client{
		Transport: tr,
	}

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	start := time.Now()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			_, _ = httpCli.Do(req.WithContext(ToContext(req.Context(), testRateLimiterKey)))
			wg.Done()
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	if elapsed < time.Duration(2*time.Second) {
		t.Fatalf("the execution time %v should be more than 2 seconds with 2 RPS of 5 requests. rate limiter is not working properly", elapsed)
	}
}
