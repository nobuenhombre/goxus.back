package domainapp

import (
	"goxus/src/internal/app/goxus/cli"
	configapp "goxus/src/internal/app/goxus/config"
)

type DomainService interface {
	Run() error
	GetConfig() *configapp.Config
}

type AppDomain struct {
	Cli    *cli.Config
	Config *configapp.Config
}

func New(cliConfig cli.Service, appConfig configapp.Service) DomainService {
	return &AppDomain{
		Cli:    cliConfig.(*cli.Config),
		Config: appConfig.Get(),
	}
}

func (d *AppDomain) Run() error {
	// Add your domain logic here
	return nil
}

func (d *AppDomain) GetConfig() *configapp.Config {
	return d.Config
}
