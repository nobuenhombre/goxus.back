package cli

import (
	"github.com/nobuenhombre/suikat/pkg/clivar"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

// Run type constants.
const (
	RunTypeInit    = "init"
	RunTypeService = "service"
)

type Service interface {
}

// Config represents the command-line interface configuration structure.
type Config struct {
	RunType string `cli:"runtype[Run type (init/service)]:string=init"`
	Config  string `cli:"config[Path to YAML config]:string=config.yaml"`
	LogFile string `cli:"log[Path to log file]:string="`
}

// New creates a new Config instance by loading values from command-line arguments.
func New() (Service, error) {
	cfg := &Config{}

	err := clivar.Load(cfg)
	if err != nil {
		return nil, ge.Pin(err)
	}

	return cfg, nil
}
