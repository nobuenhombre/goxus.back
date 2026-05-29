package configcron

import (
	configexample "goxus/src/internal/app/goxus/cron-job/jobs/example/config"
)

type CronConfig struct {
	ExampleJob configexample.ExampleJobConfig `yaml:"example_job"`
}
