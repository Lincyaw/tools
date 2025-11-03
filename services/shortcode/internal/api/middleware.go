package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// loggerMiddleware logger middleware
func loggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		statusColor := param.StatusCodeColor()
		resetColor := param.ResetColor()

		return fmt.Sprintf("%s | %s | %s | %s%3d%s | %13v | %15s | %s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			statusColor,
			param.StatusCode,
			resetColor,
			param.Latency,
			param.ClientIP,
			param.ErrorMessage,
		)
	})
}

// corsMiddleware CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimiter in-memory based simple rate limiter
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // number of requests allowed per time window
	window   time.Duration // time window
}

type visitor struct {
	requests  int
	resetTime time.Time
}

// NewRateLimiter create rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	// start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

// Allow checks if the request is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[ip]

	if !exists || now.After(v.resetTime) {
		rl.visitors[ip] = &visitor{
			requests:  1,
			resetTime: now.Add(rl.window),
		}
		return true
	}

	if v.requests < rl.rate {
		v.requests++
		return true
	}

	return false
}

// cleanupVisitors periodically cleans up expired visitor records
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.After(v.resetTime) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// rateLimitMiddleware rate limiting middleware
func rateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// timeoutMiddleware request timeout middleware
func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "request_timeout",
				"message": "Request processing timeout",
			})
			c.Abort()
		case p := <-panicChan:
			panic(p)
		case <-finished:
			return
		}
	}
}

// errorHandlerMiddleware unified error handling middleware
func errorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "internal_server_error",
					"message": "An unexpected error occurred",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}

// requestIDMiddleware request ID middleware
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		c.Set("RequestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// securityHeadersMiddleware security headers middleware
func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
