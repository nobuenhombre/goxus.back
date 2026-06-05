package v1

import (
	"github.com/gin-gonic/gin"

	"goxus/src/internal/app/goxus/api/server/router/v1/handlers"
	"goxus/src/internal/app/goxus/api/server/router/v1/middlewares"
	domainapp "goxus/src/internal/app/goxus/domain"
	"goxus/src/internal/pkg/services/ratelimit"
)

// SetupRoutes registers all v1 API routes on the given RouterGroup.
// All routes are mounted under /api/v1.
func SetupRoutes(api *gin.RouterGroup, dom domainapp.DomainService, rl ratelimit.Service) {
	v1 := api.Group("/v1")
	{
		m := middlewares.NewHttpMiddleware(dom, rl)
		authMiddleware := m.AuthTokenMiddleware()
		loginRateLimitMiddleware := m.LoginRateLimitMiddleware()

		// Handlers
		h := handlers.NewHttpHandler(dom)

		// Routes
		v1.GET("/", h.Welcome)
		v1.GET("/health", h.Health)

		// Auth
		auth := v1.Group("/auth")
		{
			auth.POST("/login", loginRateLimitMiddleware, h.LoginHandler)
		}

		// Authenticated user routes
		user := v1.Group("/user")
		user.Use(authMiddleware)
		{
			user.POST("/logout", h.LogoutHandler)
		}

		// User management
		entity := v1.Group("/entity")
		entity.Use(authMiddleware)
		{
			users := entity.Group("/user")
			{
				users.POST("/", h.CreateUser)
				users.GET("/", h.ListUsers)
				users.GET("/:id", h.GetUser)
				users.PUT("/:id", h.UpdateUser)
				users.DELETE("/:id", h.DeleteUser)

				// Restore user
				users.POST("/:id/restore", h.RestoreUser)

				// Change password
				users.PUT("/:id/password", h.ChangeUserPassword)

				// User roles
				users.GET("/:id/roles", h.GetUserRoles)
				users.POST("/:id/roles", h.AssignUserRole)
				users.DELETE("/:id/roles/:slug", h.RevokeUserRole)
			}
		}
	}
}
