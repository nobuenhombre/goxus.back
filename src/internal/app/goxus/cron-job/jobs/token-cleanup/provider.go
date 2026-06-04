package tokencleanup

import (
	"log"

	"github.com/google/wire"

	configapp "goxus/src/internal/app/goxus/config"
	domainapp "goxus/src/internal/app/goxus/domain"
)

// ProviderSet exports Wire providers for the token-cleanup cron job.
var ProviderSet = wire.NewSet(
	ProvideTokenCleanupJob,
)

// ProvideTokenCleanupJob builds the token-cleanup cron job instance (without scheduler).
// The scheduler is wired centrally in the cron-job/jobs package.
func ProvideTokenCleanupJob(appConfig configapp.Service, dom domainapp.DomainService) (*Job, func(), error) {
	cleanup := func() {
		log.Println("Token-cleanup job cleanup")
	}

	cfg := appConfig.Get()

	job, err := New(dom, &cfg.Token)
	if err != nil {
		return nil, cleanup, err
	}

	return job, cleanup, nil
}
