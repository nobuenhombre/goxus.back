package domainapp

import (
	"log"

	"github.com/google/wire"
	"goxus/src/internal/app/goxus/cli"
	configapp "goxus/src/internal/app/goxus/config"
)

// ProviderSet exports Wire providers for the domainapp package.
var ProviderSet = wire.NewSet(
	ProvideDomain,
)

// ProvideDomain creates the domain service (business-logic orchestrator).
func ProvideDomain(cliConfig cli.Service, appConfig configapp.Service) (DomainService, func(), error) {
	cleanup := func() {
		log.Println("Domain cleanup")
	}

	dom := New(cliConfig, appConfig)

	return dom, cleanup, nil
}
