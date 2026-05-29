package configcron

import (
	"github.com/google/wire"
)

// ProviderSet exports Wire providers for the cron-job config package.
var ProviderSet = wire.NewSet(
	ProvideCronConfig,
)

// ProvideCronConfig creates the cron configuration.
func ProvideCronConfig() (CronConfig, func(), error) {
	cleanup := func() {
	}

	return CronConfig{}, cleanup, nil
}
