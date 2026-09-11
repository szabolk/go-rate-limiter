package ratelimiter

import "time"

// RateLimiter decides whether a request should be allowed right now.
// Since I'm implementing multiple algorithms for testing and comparison purposes
// all algorithms will implement the Allow() function for ease of testing
type RateLimiter interface {
	Allow() bool
}

// Clock returns the current time
type Clock func() time.Time