package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	v1middlewares "goxus/src/internal/app/goxus/api/server/router/v1/middlewares"
	userdomain "goxus/src/internal/app/goxus/domain/user"
	"goxus/src/internal/pkg/db/goxus"
)

// ---- DTOs ----

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

// LoginResponse is the response body for user login.
type LoginResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// ---- Auth HTTP handlers ----

// LoginHandler handles POST /api/v1/auth/login
func (h *HttpHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.Domain.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, userdomain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, userdomain.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "user logged in successfully",
		"data": LoginResponse{
			Token:  token.Token,
			UserID: user.ID,
			Name:   user.Name,
			Email:  user.Email,
		},
	})
}

// LogoutHandler handles POST /api/v1/user/logout
// Requires a valid Bearer token via Authorization header (set by AuthTokenMiddleware).
func (h *HttpHandler) LogoutHandler(c *gin.Context) {
	tokenVal, exists := c.Get(v1middlewares.Token)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	usersToken, ok := tokenVal.(*goxus.UsersToken)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token type"})
		return
	}

	err := h.Domain.Logout(c.Request.Context(), usersToken.Token)
	if err != nil {
		if errors.Is(err, userdomain.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "user logged out successfully",
	})
}
