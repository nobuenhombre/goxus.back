package tokencleanup

import (
	"context"
	configcleantokens "goxus/src/internal/app/goxus/cron-job/jobs/token-cleanup/config"
	"log"

	"github.com/nobuenhombre/suikat/pkg/ge"

	domainapp "goxus/src/internal/app/goxus/domain"
)

// Job implements cron.Job for the expired-token cleanup task.
type Job struct {
	dom     domainapp.DomainService
	ttlDays int
}

// New creates a new token-cleanup cron job.
// ttlDays must be > 0 — the caller (config.New) ensures the default.
func New(dom domainapp.DomainService, cfg *configcleantokens.TokenCleanupConfig) (*Job, error) {
	if dom == nil {
		return nil, ge.Pin(&ge.ServiceRequiredError{
			ServiceName: "domainapp.DomainService",
		})
	}

	cfg.SetDefaults()

	return &Job{
		dom:     dom,
		ttlDays: cfg.TTLDays,
	}, nil
}

// Run executes the job. Called by cron scheduler.
func (j *Job) Run() {
	err := j.dom.DeleteExpiredTokens(context.Background(), j.ttlDays)
	if err != nil {
		log.Printf("[token-cleanup] ERROR: %v", err)
		return
	}

	log.Printf("[token-cleanup] Delete expired tokens")
}
