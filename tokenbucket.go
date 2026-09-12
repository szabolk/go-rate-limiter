package ratelimiter

import (
	"sync"
	"time"
)

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
