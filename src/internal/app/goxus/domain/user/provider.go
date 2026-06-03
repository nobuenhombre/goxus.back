package userdomain

import (
	"log"

	"goxus/src/internal/pkg/db/goxus"
	"goxus/src/internal/pkg/services/rbac"

	"github.com/google/wire"
)

// ProviderSet exports Wire providers for the userdomain package.
var ProviderSet = wire.NewSet(
	ProvideUserService,
)

// ProvideUserService creates the user domain service with authorization.
// It builds a pure business-logic service, then wraps it with an
// authorization decorator that enforces RBAC permission checks.
func ProvideUserService(dbRepo *goxus.DbGoxusRepo, rbacSvc rbac.Service) (Service, func(), error) {
	cleanup := func() {
		log.Println("User domain service cleanup")
	}

	// Pure business logic layer
	raw := New(dbRepo, rbacSvc)

	// Authorization decorator layer
	svc := NewAuthorized(raw, rbacSvc)

	return svc, cleanup, nil
}
