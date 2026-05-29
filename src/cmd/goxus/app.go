package main

import (
	"log"

	"github.com/robfig/cron/v3"
	"goxus/src/internal/app/goxus/api/server"
	"goxus/src/internal/app/goxus/cli"
	domainapp "goxus/src/internal/app/goxus/domain"
)

// IApp is the top-level application orchestrator interface.
type IApp interface {
	Run() error
}

// App is the top-level application orchestrator.
type App struct {
	cliConfig  cli.Service
	dom        domainapp.DomainService
	httpServer server.IHTTPServer
	cronJob    *cron.Cron
}

// Run executes the application based on CLI configuration.
func (a *App) Run() error {
	cliConfig := a.cliConfig.(*cli.Config)

	switch cliConfig.RunType {
	case cli.RunTypeInit:
		log.Printf("Running in init mode")

	case cli.RunTypeService:
		if a.cronJob != nil {
			a.cronJob.Start()
		}

		return a.httpServer.Run()

	default:
		log.Printf("Running in default mode")
		return a.dom.Run()
	}

	return nil
}

// newApp is the Wire provider for the top-level application.
func newApp(cliConfig cli.Service, dom domainapp.DomainService, httpServer server.IHTTPServer, cronJob *cron.Cron) (IApp, func(), error) {
	cleanup := func() {
		log.Println("App cleanup")
		if cronJob != nil {
			cronJob.Stop()
		}
	}

	return &App{
		cliConfig:  cliConfig,
		dom:        dom,
		httpServer: httpServer,
		cronJob:    cronJob,
	}, cleanup, nil
}
