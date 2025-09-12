package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type RateLimiter struct {
	requests    map[string]*WindowData
	mutex       sync.RWMutex
	permitLimit int
	window      time.Duration
	queueLimit  int
}

type WindowData struct {
	count     int
	queue     []time.Time
	windowEnd time.Time
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(permitLimit int, window time.Duration, queueLimit int) *RateLimiter {
	rl := &RateLimiter{
		requests:    make(map[string]*WindowData),
		permitLimit: permitLimit,
		window:      window,
		queueLimit:  queueLimit,
	}

	// Cleanup goroutine to remove expired entries
	go rl.cleanup()

	return rl
}

// cleanup removes expired entries from the rate limiter
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for key, data := range rl.requests {
			if now.After(data.windowEnd) {
				delete(rl.requests, key)
			}
		}
		rl.mutex.Unlock()
	}
}

// isAllowed checks if the request is allowed for the given identifier
func (rl *RateLimiter) isAllowed(identifier string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Get or create window data for this identifier
	data, exists := rl.requests[identifier]
	if !exists || now.After(data.windowEnd) {
		// Create new window
		rl.requests[identifier] = &WindowData{
			count:     1,
			queue:     []time.Time{now},
			windowEnd: now.Add(rl.window),
		}
		return true
	}

	// Clean up expired requests from queue
	validRequests := make([]time.Time, 0, len(data.queue))
	for _, reqTime := range data.queue {
		if now.Sub(reqTime) < rl.window {
			validRequests = append(validRequests, reqTime)
		}
	}
	data.queue = validRequests
	data.count = len(validRequests)

	// Check if under limit
	if data.count < rl.permitLimit {
		data.queue = append(data.queue, now)
		data.count++
		return true
	}

	// Check queue limit
	if len(data.queue) >= rl.queueLimit {
		return false
	}

	return false
}

// Middleware returns Echo middleware function for rate limiting
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Use client IP as identifier (you can customize this)
			identifier := c.RealIP()

			if !rl.isAllowed(identifier) {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}

			return next(c)
		}
	}
}

// MiddlewareWithKeyFunc returns Echo middleware with custom key function
func (rl *RateLimiter) MiddlewareWithKeyFunc(keyFunc func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			identifier := keyFunc(c)

			if !rl.isAllowed(identifier) {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}

			return next(c)
		}
	}
}

// GetStats returns current statistics for an identifier
func (rl *RateLimiter) GetStats(identifier string) (count int, remaining int, resetTime time.Time) {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	data, exists := rl.requests[identifier]
	if !exists {
		return 0, rl.permitLimit, time.Now().Add(rl.window)
	}

	now := time.Now()
	if now.After(data.windowEnd) {
		return 0, rl.permitLimit, now.Add(rl.window)
	}

	// Clean up expired requests for accurate count
	validCount := 0
	for _, reqTime := range data.queue {
		if now.Sub(reqTime) < rl.window {
			validCount++
		}
	}

	remaining = rl.permitLimit - validCount
	if remaining < 0 {
		remaining = 0
	}

	return validCount, remaining, data.windowEnd
}
