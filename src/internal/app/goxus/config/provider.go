package configapp

import (
	"log"

	"goxus/src/internal/app/goxus/cli"

	"github.com/google/wire"
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

// ProviderSet exports Wire providers for the configapp package.
var ProviderSet = wire.NewSet(
	ProvideConfigApp,
	ProvideDBConfig,
)

// ProvideConfigApp loads the YAML configuration from the CLI config file path.
func ProvideConfigApp(cliConfig cli.Service) (Service, func(), error) {
	cleanup := func() {
		log.Println("App config cleanup")
	}

	cfg, err := New(cliConfig.(*cli.Config).Config)
	if err != nil {
		return nil, cleanup, ge.Pin(err)
	}

	return cfg, cleanup, nil
}

// ProvideDBConfig provides the database connection config from the app config.
func ProvideDBConfig(appConfig Service) (*pgxdb.Config, error) {
	return new(appConfig.Get().DB), nil
}
