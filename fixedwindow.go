package ratelimiter

import (
	"time"
	"sync"
)

// FixedWindow is a rate limiter that allows a certain number of requests in a fixed time window.
// Something to do after: Once finished, drop the mutex and use atomic operations through Redis or something similar to avoid locking and unlocking the mutex for every request. 
// something something performance or what not
type FixedWindow struct {
	mu          sync.Mutex
	limit       int
	window      time.Duration
	clock       Clock
	windowStart time.Time
	count       int
}