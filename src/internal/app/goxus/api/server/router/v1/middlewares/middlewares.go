package middlewares

import (
	"goxus/src/internal/pkg/services/ratelimit"

	domainapp "goxus/src/internal/app/goxus/domain"
)

// HttpMiddleware holds HTTP middleware methods.
type HttpMiddleware struct {
	Domain      domainapp.DomainService
	RateLimiter ratelimit.Service
}

// NewHttpMiddleware creates a new HttpMiddleware.
func NewHttpMiddleware(dom domainapp.DomainService, rl ratelimit.Service) (mid *HttpMiddleware) {
	mid = new(HttpMiddleware)
	mid.Domain = dom
	mid.RateLimiter = rl

	return mid
}
