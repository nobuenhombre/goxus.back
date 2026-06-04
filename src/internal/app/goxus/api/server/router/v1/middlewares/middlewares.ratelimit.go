package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LoginRateLimitMiddleware limits the number of login attempts per client IP.
// It uses the configured in-memory sliding-window rate limiter.
// Returns 429 Too Many Requests with Retry-After header when rate-limited.
func (mid *HttpMiddleware) LoginRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// identify client by IP
		key := c.ClientIP()

		if !mid.RateLimiter.Allow(key) {
			remaining := mid.RateLimiter.ResetAfter(key)
			c.Header("Retry-After", strconv.Itoa(int(remaining.Seconds())))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many login attempts",
				"retry_after": remaining.Seconds(),
			})
			return
		}

		c.Next()
	}
}
