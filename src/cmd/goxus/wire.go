//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"goxus/src/internal/app/goxus/api/server"
	"goxus/src/internal/app/goxus/cli"
	configapp "goxus/src/internal/app/goxus/config"
	examplejobs "goxus/src/internal/app/goxus/cron-job/jobs/example"
	domainapp "goxus/src/internal/app/goxus/domain"
	logfile "goxus/src/internal/app/goxus/log"
)

// initializeApp is the Wire injector entrypoint. It aggregates all ProviderSets
// and constructs the top-level application. No logic belongs here.
func initializeApp() (IApp, func(), error) {
	wire.Build(
		cli.ProviderSet,
		logfile.ProviderSet,
		configapp.ProviderSet,
		domainapp.ProviderSet,
		examplejobs.ProviderSet,
		server.ProviderSet,
		newApp,
	)
	return nil, nil, nil
}
