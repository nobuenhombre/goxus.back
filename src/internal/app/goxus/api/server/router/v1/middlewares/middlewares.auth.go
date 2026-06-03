package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	userdomain "goxus/src/internal/app/goxus/domain/user"
)

const (
	Token = "token"
	User  = "user"
)

// AuthTokenMiddleware validates the Bearer token from Authorization header
// and sets the token and user in the gin context.
func (mid *HttpMiddleware) AuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")
		tokenIdentity := strings.TrimLeft(tokenHeader, "Bearer")
		tokenIdentity = strings.TrimSpace(tokenIdentity)

		if tokenIdentity == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
			return
		}

		user, token, err := mid.Domain.ValidateToken(c, tokenIdentity)
		if err != nil {
			if errors.Is(err, userdomain.ErrTokenNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user token not found"})
				return
			}

			if errors.Is(err, userdomain.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}

			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}

		c.Set(Token, token)
		c.Set(User, user)

		c.Next()
	}
}
