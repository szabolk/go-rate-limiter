package ratelimiter

import (
	"sync"
	"time"
)

// FixedWindow is a rate limiter that allows a certain number of requests in a fixed time window.
// Something to do after: Once finished, drop the mutex and use atomic operations through Redis or something similar to avoid locking and unlocking the mutex for every request.
// something something performance or what not
// Note to self: Issue with fixed window is that a user can do X number of requests at xx:59, max out the limit, and then one minute later at xx:00, max out the limit again
// Ex. 100 requests an hour, user does 100 at xx:59, time becomes xx:00 -> user request count back to 0 -> user can now request 100 more times
// Pretty much able to do double the # of requests an hour in like 2 minutes
type FixedWindow struct {
	mu          sync.Mutex
	limit       int
	window      time.Duration
	clock       Clock
	windowStart time.Time
	count       int
}

func NewFixedWindow(clock Clock, limit int, window time.Duration) *FixedWindow {
	return &FixedWindow{
		limit:       limit,
		window:      window,
		clock:       clock,
		windowStart: clock(),
	}
}

// Allow checks if a request is allowed under the rate limit
func (f *FixedWindow) Allow() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	currentTime := f.clock()

	if currentTime.Sub(f.windowStart) >= f.window { //Check if window is expired
		f.windowStart = currentTime
		f.count = 0
	}
	if f.count < f.limit {
		f.count++
		return true
	}
	return false
}
