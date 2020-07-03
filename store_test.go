package ratelimiter

import "testing"

func TestStore(t *testing.T) {
	s := NewStore()
	s.Register("mock", 10, 10)
	// error out on existing key
	if err := s.Register("mock", 1, 1); err == nil {
		t.Fatalf("expected error but got nil")
	}
	limiter, err := s.Get("mock")
	if err != nil {
		t.Fatal(err)
	}
	if limiter == nil {
		t.Fatalf("expected ratelimiter instance but got nil")
	}
	// error out on key does not exist
	if _, err := s.Get("mock2"); err == nil {
		t.Fatalf("expected error but got nil")
	}
}
