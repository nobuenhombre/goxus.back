package rbac

import (
	"log"

	"goxus/src/internal/pkg/db/goxus"

	"github.com/google/wire"
)

// ProviderSet exports Wire providers for the rbac package.
var ProviderSet = wire.NewSet(
	ProvideRbac,
)

// ProvideRbac creates the RBAC service.
func ProvideRbac(dbRepo *goxus.DbGoxusRepo) (Service, func(), error) {
	cleanup := func() {
		log.Println("RBAC service cleanup")
	}

	svc := New(dbRepo)

	return svc, cleanup, nil
}
