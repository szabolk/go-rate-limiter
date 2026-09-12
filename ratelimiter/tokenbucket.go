package ratelimiter

import (
	"sync"
	"time"
)

// Notes for self: Basically while a bucket has tokens, requests are accepted
type TokenBucket struct {
	mu         sync.Mutex
	capacity   float64
	refillRate float64
	clock      Clock
	tokens     float64
	lastRefill time.Time
}

func NewTokenBucket(capacity float64, refillRate float64, clock Clock) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		refillRate: refillRate,
		clock:      clock,
		tokens:     capacity,
		lastRefill: clock(),
	}
}

func (t *TokenBucket) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	currentTime := t.clock()
	elapsedTime := currentTime.Sub(t.lastRefill).Seconds()
	t.tokens += elapsedTime * t.refillRate

	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}

	t.lastRefill = currentTime
	if t.tokens < 1 {
		return false
	}
	t.tokens--
	return true

}
