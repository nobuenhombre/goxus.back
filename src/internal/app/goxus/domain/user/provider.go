package userdomain

import (
	"log"

	configapp "goxus/src/internal/app/goxus/config"
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
func ProvideUserService(dbRepo *goxus.DbGoxusRepo, rbacSvc rbac.Service, appConfig configapp.Service) (Service, func(), error) {
	cleanup := func() {
		log.Println("User domain service cleanup")
	}

	cfg := appConfig.Get()
	cfg.Storage.SetDefaults()

	// Pure business logic layer
	raw := New(dbRepo, rbacSvc, cfg)

	// Authorization decorator layer
	svc := NewAuthorized(raw, rbacSvc)

	return svc, cleanup, nil
}
