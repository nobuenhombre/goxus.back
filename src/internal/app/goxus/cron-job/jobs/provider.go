package cronjobs

import (
	"log"

	"github.com/google/wire"
	"github.com/robfig/cron/v3"

	configapp "goxus/src/internal/app/goxus/config"
	examplejobs "goxus/src/internal/app/goxus/cron-job/jobs/example"
	tokencleanup "goxus/src/internal/app/goxus/cron-job/jobs/token-cleanup"
)

// ProviderSet exports Wire providers for the central cron scheduler.
// It aggregates all job instances and builds a single *cron.Cron.
var ProviderSet = wire.NewSet(
	examplejobs.ProviderSet,
	tokencleanup.ProviderSet,
	ProvideCronScheduler,
)

// ProvideCronScheduler builds the single cron scheduler with all jobs registered.
// Each job is registered only if enabled in config.
func ProvideCronScheduler(
	exampleJob *examplejobs.Job,
	tokenCleanupJob *tokencleanup.Job,
	appConfig configapp.Service,
) (*cron.Cron, func(), error) {
	cfg := appConfig.Get()

	c := cron.New()

	if cfg.Cron.ExampleJob.Enabled {
		_, err := c.AddJob(cfg.Cron.ExampleJob.Schedule, exampleJob)
		if err != nil {
			return nil, nil, err
		}
		log.Println("[cron] Registered example_job with schedule:", cfg.Cron.ExampleJob.Schedule)
	}

	if cfg.Cron.TokenCleanupJob.Enabled {
		_, err := c.AddJob(cfg.Cron.TokenCleanupJob.Schedule, tokenCleanupJob)
		if err != nil {
			return nil, nil, err
		}
		log.Println("[cron] Registered token_cleanup_job with schedule:", cfg.Cron.TokenCleanupJob.Schedule)
	}

	cleanup := func() {
		log.Println("[cron] Stopping scheduler")
		c.Stop()
	}

	return c, cleanup, nil
}
