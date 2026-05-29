package configexample

type ExampleJobConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Schedule string `yaml:"schedule,omitempty"`
}
