package router

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"goxus/src/internal/app/goxus/api/server/router/handlers"
	"goxus/src/internal/app/goxus/api/server/router/middlewares"
	"goxus/src/internal/app/goxus/api/server/router/v1"
	domainapp "goxus/src/internal/app/goxus/domain"
	"goxus/src/internal/app/goxus/version"
)

// HTTPRouter wraps the Gin engine with versioned API routes.
type HTTPRouter struct {
	Router      *gin.Engine
	Handlers    *handlers.HttpHandler
	Middlewares *middlewares.HttpMiddleware
}

// NewHTTPRouter creates a new HTTPRouter.
func NewHTTPRouter(logFile *os.File, dom domainapp.DomainService) (router *HTTPRouter) {
	router = new(HTTPRouter)

	if logFile != nil {
		gin.DisableConsoleColor()
		gin.DefaultWriter = io.MultiWriter(logFile)
	}

	router.Handlers = handlers.NewHttpHandler(dom)
	router.Middlewares = middlewares.NewHttpMiddleware(dom)

	router.Router = gin.Default()
	router.Router.Use(router.Middlewares.CORSMiddleware())
	router.Router.Use(router.Middlewares.APILoggerMiddleware())

	router.SetupRoutes(dom)

	return
}

// SetupRoutes registers all HTTP routes including versioned API groups.
func (r *HTTPRouter) SetupRoutes(dom domainapp.DomainService) {
	// Unversioned health check
	r.Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"app_version": version.Version,
			"server_time": time.Now().Format(time.RFC3339),
		})
	})

	// Versioned API group
	api := r.Router.Group("/api")
	{
		v1.SetupRoutes(api, dom)
	}
}
