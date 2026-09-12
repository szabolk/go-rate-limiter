package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	rl "ratelimiter/ratelimiter"
)

func runDemo(name string, limiter rl.RateLimiter, numWindows, goroutines, requestsPerGoroutine int, window time.Duration) {
	var totalAllowed int64
	total := int64(numWindows * goroutines * requestsPerGoroutine)

	for w := range numWindows {
		var wg sync.WaitGroup
		var windowAllowed int64

		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < requestsPerGoroutine; j++ {
					if limiter.Allow() {
						atomic.AddInt64(&windowAllowed, 1)
					}
				}
			}()
		}
		wg.Wait()

		atomic.AddInt64(&totalAllowed, windowAllowed)
		fmt.Printf("%-15s window %d: %d/%d allowed\n", name, w+1, windowAllowed, int64(goroutines*requestsPerGoroutine))

		time.Sleep(window + 50*time.Millisecond) // let the window fully roll over
	}

	fmt.Printf("%-15s TOTAL: %d/%d requests allowed\n\n", name, totalAllowed, total)
}

func main() {
	const limit = 20
	const window = time.Second
	const goroutines = 10
	const perGoroutine = 12 // 120 reqs total

	fmt.Printf("Limit: %d per %v | %d goroutines x %d requests each (%d total)\n\n",
		limit, window, goroutines, perGoroutine, goroutines*perGoroutine)
	const numWindows = 10

	// Fixed window should show results of limit / (goroutines*perGoroutine) for all windows
	runDemo("FixedWindow", rl.NewFixedWindow(time.Now, limit, window), numWindows, goroutines, perGoroutine, window)
	// Token bucket will depend on the refill rate mainly
	runDemo("TokenBucket", rl.NewTokenBucket(float64(limit), float64(limit)/2, time.Now), numWindows, goroutines, perGoroutine, window)
	// Sliding window depends on how long the system sleeps for
	runDemo("SlidingWindow", rl.NewSlidingWindow(time.Now, limit, window), numWindows, goroutines, perGoroutine, window)
}
