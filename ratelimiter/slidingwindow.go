package ratelimiter

import (
	"sync"
	"time"
)

// Notes for myself: Sliding window, while more complicated, is better than fixed window in a few ways
// Unlike fixed window, SW doesn't suffer from the all requests at the end of the time period and then immediately
// request more at the start of the next one.
type SlidingWindow struct {
	mu        sync.Mutex
	limit     int
	window    time.Duration
	clock     Clock
	currStart time.Time
	currCount int
	prevCount int
}

func NewSlidingWindow(clock Clock, limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		limit:     limit,
		window:    window,
		clock:     clock,
		currStart: clock(),
	}
}

func (s *SlidingWindow) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTime := s.clock()
	elapsedTime := currentTime.Sub(s.currStart)

	if elapsedTime >= s.window {
		numWindows := int(elapsedTime / s.window) // Calculates how many windows have passed based on window duration
		if numWindows == 1 {                      // Base case: only one window passed by (at beginning)
			s.prevCount = s.currCount
		} else { // More than one window is passed => logically previous window no longer previous window => its count is 0
			s.prevCount = 0
		}
		s.currCount = 0
		s.currStart = s.currStart.Add(time.Duration(numWindows) * s.window) // Alter currStart by calculating how much time has elapsed
		elapsedTime = currentTime.Sub(s.currStart)
	}

	percentOverlap := float64(s.window-elapsedTime) / float64(s.window)
	if percentOverlap < 0 {
		percentOverlap = 0
	}
	weightedCount := float64(s.prevCount)*percentOverlap + float64(s.currCount) // Estimates # of requests by using the overlap percent and calculating the count

	if weightedCount >= float64(s.limit) {
		return false
	}
	s.currCount++
	return true
}
