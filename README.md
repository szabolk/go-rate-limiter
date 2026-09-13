# Rate Limiter in Go
This was my first exposure to the idea of a rate limiter (as a developer). I had an idea of what a rate limiter did but I never really considered how one worked and what they used in terms of algorithm to allow/deny requests sent by users. This was also my introduction to Go, so I wanted to do something that had some semblance of real world relevance.

## Overview
This project implements three common rate limiting algorithms, all sharing a single `RateLimiter` interface with one method: `Allow() bool`. `Allow()` essentially just determines if a request is allowed to go through or not depending on the defined request amount limits, etc.
 
Running `main.go` runs all three limiters side by side under the same concurrent load (multiple goroutines firing requests across several time windows). Just a way to show how the different algorithms work depending on constraints.

## Algorithms
 
### Fixed Window
 
Counts requests in a fixed-size time window (Ex. 20 requests per second or whatever request per timeframe is chosen). Once the window expires, the count is reset to 0 and the user can request again.
It was the simplest to implement, but suffers from the case when a user requests at the end of the window and then requests at the beginning of the next. Ex. 20 requests per minute. User does 20 requests at the 59th second, and then once the next minute arrives and resets the request count, the user can do 20 requests again. Ultimately, this means that the user effective did 40 requests (double what is permitted in this case) in less than a minute.
 
### Sliding Window
 
Improves on the fixed window by looking at a weighted combination of the current window's count and the previous window's count, based on how much the current window overlaps with the previous one. Gets rid of the problem that fixed window has by determining how many requests a user has done at all times within a certain time frame. 

Ex. Say that during the previous window and current window have a 50% overlap. Let's say the previous window request count is 10 and the current window request count is 20. Using the formula:

$$\Huge \text{weightedCount} = \text{prevCount} \times \text{percentOverlap} + \text{currCount}$$

We get: 

$$\Huge \text{weightedCount} = \text{10} \times \text{0.5} + \text{20}$$

$$\Huge \text{weightedCount} = \text{25}$$

Essentially, this means that out of the 20 requests the user submitted in the current window, only 15 of those will be accepted, with the other 5 being declined. For the user to be able to send in more requests, they would have to wait for overlap of the previous window to decrease as the window moves along.

### Token Bucket
 
Maintains a bucket of tokens that refills continuously over time up to a maximum capacity. Each allowed request consumes one token; if the bucket is empty, the request is denied. The number of requests a user is able to do in a given time frame is determined by the bucket's refill rate.

## Time and Testing

The `Clock` type (`func() time.Time`) is injected into each limiter rather than calling `time.Now()` directly. I originally tried `time.Now()`, but having to wait in real time for testing got annoying real fast. The `Clock` type allowed me to test things instantly.

## Running
```
go run main.go
```
This runs all three algorithms against the same configuration (check `main.go`). Changing the values in `main.go`, such as `limit`, `window`, `goroutines`, and `perGorountine` will have an effect on the algorithms.

## Requirements
Check `go.mod`.

## Future Endeavours
- Test files that show the intricacies of each algorithm.
- Instead of using mutex locks to ensure that concurrent requests don't cause race conditions, use atomic operations through Redis.
- TBD
