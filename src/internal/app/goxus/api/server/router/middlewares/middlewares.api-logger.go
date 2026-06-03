package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// APILoggerMiddleware logs each request with method, path, status, and latency.
func (mid *HttpMiddleware) APILoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		log.Printf("[v1] %s %s %d %v", c.Request.Method, path, c.Writer.Status(), latency)
	}
}
