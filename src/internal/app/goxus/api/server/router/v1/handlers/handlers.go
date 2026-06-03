package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	domainapp "goxus/src/internal/app/goxus/domain"
	"goxus/src/internal/app/goxus/version"
)

// HttpHandler holds v1-specific HTTP handler methods.
type HttpHandler struct {
	Domain domainapp.DomainService
}

// NewHttpHandler creates a new v1 HttpHandler.
func NewHttpHandler(dom domainapp.DomainService) *HttpHandler {
	return &HttpHandler{
		Domain: dom,
	}
}

// Welcome returns a welcome message for the v1 API.
func (h *HttpHandler) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "v1",
		"message": "Welcome to goxus API v1",
	})
}

// Health returns a simple health check response.
func (h *HttpHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"app_version": version.Version,
		"server_time": time.Now().Format(time.RFC3339),
	})
}
