package examplejobs

import (
	"log"

	"github.com/google/wire"
	"github.com/nobuenhombre/suikat/pkg/ge"
	"github.com/robfig/cron/v3"
	domainapp "goxus/src/internal/app/goxus/domain"
)

// ProviderSet exports Wire providers for the example cron job.
var ProviderSet = wire.NewSet(
	ProvideExampleJobs,
)

// ProvideExampleJobs builds the cron scheduler with the example job.
func ProvideExampleJobs(dom domainapp.DomainService) (*cron.Cron, func(), error) {
	cleanup := func() {
		log.Println("Example jobs cleanup")
	}

	cfg := dom.GetConfig()

	exampleJob, err := New(dom)
	if err != nil {
		return nil, cleanup, ge.Pin(err)
	}

	c := cron.New()

	if cfg.Cron.ExampleJob.Enabled {
		_, err = c.AddJob(cfg.Cron.ExampleJob.Schedule, exampleJob)
		if err != nil {
			return nil, cleanup, ge.Pin(err)
		}
	}

	return c, cleanup, nil
}
