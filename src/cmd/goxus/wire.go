//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"goxus/src/internal/app/goxus/api/server"
	"goxus/src/internal/app/goxus/cli"
	configapp "goxus/src/internal/app/goxus/config"
	cronjobs "goxus/src/internal/app/goxus/cron-job/jobs"
	domainapp "goxus/src/internal/app/goxus/domain"
	settingsdomain "goxus/src/internal/app/goxus/domain/settings"
	userdomain "goxus/src/internal/app/goxus/domain/user"
	logfile "goxus/src/internal/app/goxus/log"
	"goxus/src/internal/pkg/db/goxus"
	"goxus/src/internal/pkg/services/ratelimit"
	"goxus/src/internal/pkg/services/rbac"
)

// initializeApp is the Wire injector entrypoint. It aggregates all ProviderSets
// and constructs the top-level application. No logic belongs here.
func initializeApp() (IApp, func(), error) {
	wire.Build(
		cli.ProviderSet,
		logfile.ProviderSet,
		configapp.ProviderSet,
		goxus.ProviderSet,
		rbac.ProviderSet,
		ratelimit.ProviderSet,
		userdomain.ProviderSet,
		settingsdomain.ProviderSet,
		domainapp.ProviderSet,
		cronjobs.ProviderSet,
		server.ProviderSet,
		newApp,
	)
	return nil, nil, nil
}
