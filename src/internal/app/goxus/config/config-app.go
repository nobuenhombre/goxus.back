package configapp

import (
	configserver "goxus/src/internal/app/goxus/api/server/config"
	configcron "goxus/src/internal/app/goxus/cron-job/config"
	configexample "goxus/src/internal/app/goxus/cron-job/jobs/example/config"
	configcleantokens "goxus/src/internal/app/goxus/cron-job/jobs/token-cleanup/config"
	configratelimit "goxus/src/internal/pkg/services/ratelimit/config"

	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
	"github.com/nobuenhombre/suikat/pkg/fico"
	"github.com/nobuenhombre/suikat/pkg/ge"
	"gopkg.in/yaml.v3"
)

type Service interface {
	Load(fileName string) error
	Save(fileName string) error
	Get() *Config
}

type HostsConfig struct {
	API configserver.HTTPServerConfig `yaml:"api,omitempty"`
}

type LogConfig struct {
	Quiet bool `yaml:"quiet,omitempty"`
}

type Config struct {
	DB        pgxdb.Config                         `yaml:"db,omitempty"`
	Hosts     HostsConfig                          `yaml:"hosts,omitempty"`
	Log       LogConfig                            `yaml:"log,omitempty"`
	Cron      configcron.CronConfig                `yaml:"cron,omitempty"`
	RateLimit configratelimit.LoginRateLimitConfig `yaml:"rate_limit,omitempty"`
	Example   configexample.ExampleJobConfig       `yaml:"example,omitempty"`
	Token     configcleantokens.TokenCleanupConfig `yaml:"token,omitempty"`
}

func New(fileName string) (Service, error) {
	cfg := &Config{}

	err := cfg.Load(fileName)
	if err != nil {
		return nil, ge.Pin(err)
	}

	return cfg, nil
}

func (c *Config) Load(fileName string) error {
	txtConfigFile := fico.TxtFile(fileName)

	configData, err := txtConfigFile.Read()
	if err != nil {
		return ge.Pin(err)
	}

	err = yaml.Unmarshal([]byte(configData), c)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

func (c *Config) Save(fileName string) error {
	txtConfigFile := fico.TxtFile(fileName)

	configData, err := yaml.Marshal(c)
	if err != nil {
		return ge.Pin(err)
	}

	err = txtConfigFile.Write(string(configData))
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

func (c *Config) Get() *Config {
	return c
}
