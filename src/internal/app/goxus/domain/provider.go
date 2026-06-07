package domainapp

import (
	"log"

	"goxus/src/internal/app/goxus/cli"
	configapp "goxus/src/internal/app/goxus/config"
	settingsdomain "goxus/src/internal/app/goxus/domain/settings"
	userdomain "goxus/src/internal/app/goxus/domain/user"
	"goxus/src/internal/pkg/services/rbac"

	"github.com/google/wire"
)

// ProviderSet exports Wire providers for the domainapp package.
var ProviderSet = wire.NewSet(
	ProvideDomain,
)

// ProvideDomain creates the domain service (business-logic orchestrator).
func ProvideDomain(cliConfig cli.Service, appConfig configapp.Service, rbacService rbac.Service, userService userdomain.Service, settingsService settingsdomain.Service) (DomainService, func(), error) {
	cleanup := func() {
		log.Println("Domain cleanup")
	}

	dom := New(cliConfig, appConfig, rbacService, userService, settingsService)

	return dom, cleanup, nil
}
