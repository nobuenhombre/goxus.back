package middlewares

import (
	domainapp "goxus/src/internal/app/goxus/domain"
)

// HttpMiddleware holds HTTP middleware methods.
type HttpMiddleware struct {
	Domain domainapp.DomainService
}

// NewHttpMiddleware creates a new HttpMiddleware.
func NewHttpMiddleware(dom domainapp.DomainService) (mid *HttpMiddleware) {
	mid = new(HttpMiddleware)
	mid.Domain = dom
	return mid
}
