package main

import "testing"

func TestLoginRateLimiterAllowsConfiguredBurst(t *testing.T) {
	const requestsPerMinute = 10
	store := newLoginRateLimiter(requestsPerMinute)

	for request := 1; request <= requestsPerMinute; request++ {
		allowed, err := store.Allow("127.0.0.1")
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", request, err)
		}
		if !allowed {
			t.Fatalf("request %d: expected initial burst to be allowed", request)
		}
	}

	allowed, err := store.Allow("127.0.0.1")
	if err != nil {
		t.Fatalf("request after burst: unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("request after burst: expected rate limiter to reject it")
	}
}
