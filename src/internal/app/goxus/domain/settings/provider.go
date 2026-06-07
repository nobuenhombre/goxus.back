package settingsdomain

import (
	"log"

	"goxus/src/internal/pkg/db/goxus"

	"github.com/google/wire"
)

// ProviderSet exports Wire providers for the settingsdomain package.
var ProviderSet = wire.NewSet(
	ProvideSettingsService,
)

// ProvideSettingsService creates the settings domain service.
func ProvideSettingsService(dbRepo *goxus.DbGoxusRepo) (Service, func(), error) {
	cleanup := func() {
		log.Println("Settings domain service cleanup")
	}

	svc := New(dbRepo)

	return svc, cleanup, nil
}
