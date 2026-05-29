package examplejobs

import (
	"log"

	"github.com/nobuenhombre/suikat/pkg/ge"
	"github.com/robfig/cron/v3"
	domainapp "goxus/src/internal/app/goxus/domain"
)

// Job implements cron.Job for the example scheduled task.
type Job struct {
	dom domainapp.DomainService
}

// New creates a new example cron job.
func New(dom domainapp.DomainService) (cron.Job, error) {
	if dom == nil {
		return nil, ge.Pin(&ge.ServiceRequiredError{
			ServiceName: "domainapp.DomainService",
		})
	}

	return &Job{
		dom: dom,
	}, nil
}

// Run executes the job. Called by cron scheduler.
func (j *Job) Run() {
	log.Println("Running example job")
	// Add your job logic here
}
