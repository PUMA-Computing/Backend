package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
	"log"
	"net/http"
	"sync"
)

var (
	limiterMap = make(map[string]ratelimit.Limiter)
	mu         sync.Mutex
)

// RateLimiterMiddleware creates a rate limiter for each IP address
func RateLimiterMiddleware(rate int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if _, exists := limiterMap[ip]; !exists {
			log.Printf("Creating new rate limiter for IP: %s", ip)
			limiterMap[ip] = ratelimit.New(rate)
		}
		limiter := limiterMap[ip]
		mu.Unlock()

		// Wait for the rate limiter to allow the request
		limiter.Take()

		// Check if the request should be rate limited
		if rate == 0 {
			log.Printf("Rate limit exceeded for IP: %s", ip)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		// Optionally, add a header to indicate rate limiting status
		c.Writer.Header().Set("X-Rate-Limit", "active")

		// Process the request
		c.Next()
	}
}
