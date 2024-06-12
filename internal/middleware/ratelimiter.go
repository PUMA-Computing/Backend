package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	limiterMap = make(map[string]*RateLimiter)
	mu         sync.Mutex
)

type RateLimiter struct {
	tokens         int
	lastRefill     time.Time
	refillInterval time.Duration
	maxTokens      int
}

func NewRateLimiter(maxTokens int, refillInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:         maxTokens,
		lastRefill:     time.Now(),
		refillInterval: refillInterval,
		maxTokens:      maxTokens,
	}
}

func (rl *RateLimiter) Allow() (bool, int, time.Time) {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// Refill tokens
	refillAmount := int(elapsed / rl.refillInterval)
	if refillAmount > 0 {
		rl.tokens = min(rl.tokens+refillAmount, rl.maxTokens)
		rl.lastRefill = now
	}

	// Calculate next refill time
	nextRefill := rl.lastRefill.Add(rl.refillInterval)

	// Check if we have tokens available
	if rl.tokens > 0 {
		rl.tokens--
		return true, rl.tokens, nextRefill
	}

	return false, rl.tokens, nextRefill
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func RateLimiterMiddleware(maxTokens int, refillInterval time.Duration, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiterKey := ip + ":" + key

		mu.Lock()
		if _, exists := limiterMap[limiterKey]; !exists {
			limiterMap[limiterKey] = NewRateLimiter(maxTokens, refillInterval)
		}
		limiter := limiterMap[limiterKey]
		mu.Unlock()

		allowed, tokensLeft, nextRefill := limiter.Allow()

		// Add headers to indicate rate limiting status
		c.Writer.Header().Set("X-Rate-Limit-Tokens-Left", strconv.Itoa(tokensLeft))
		c.Writer.Header().Set("X-Rate-Limit-Next-Refill", nextRefill.Format(time.RFC1123))

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"tokens_left": tokensLeft,
				"next_refill": nextRefill.Format(time.RFC1123),
			})
			return
		}

		c.Next()
	}
}
