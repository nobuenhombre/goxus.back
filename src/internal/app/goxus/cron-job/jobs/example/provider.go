package examplejobs

import (
	configapp "goxus/src/internal/app/goxus/config"
	"log"

	"github.com/google/wire"

	domainapp "goxus/src/internal/app/goxus/domain"

	"github.com/nobuenhombre/suikat/pkg/ge"
)

// ProviderSet exports Wire providers for the example job.
var ProviderSet = wire.NewSet(
	ProvideExampleJob,
)

// ProvideExampleJob builds the example cron job instance (without scheduler).
// The scheduler is wired centrally in the cron-job/jobs package.
func ProvideExampleJob(appConfig configapp.Service, dom domainapp.DomainService) (*Job, func(), error) {
	cleanup := func() {
		log.Println("Example job cleanup")
	}

	cfg := appConfig.Get()

	job, err := New(dom, &cfg.Example)
	if err != nil {
		return nil, cleanup, ge.Pin(err)
	}

	return job, cleanup, nil
}
