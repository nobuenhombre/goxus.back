package configcron

type CronJobConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Schedule string `yaml:"schedule,omitempty"`
}

type CronConfig struct {
	ExampleJob      CronJobConfig `yaml:"example_job"`
	TokenCleanupJob CronJobConfig `yaml:"token_cleanup_job,omitempty"`
}
