package server

import (
	"log"

	configapp "goxus/src/internal/app/goxus/config"
	domainapp "goxus/src/internal/app/goxus/domain"
	logfile "goxus/src/internal/app/goxus/log"

	"github.com/google/wire"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

// ProviderSet exports Wire providers for the server package.
var ProviderSet = wire.NewSet(
	ProvideAPI,
)

// ProvideAPI creates the HTTP API server (Gin) with graceful shutdown support.
func ProvideAPI(configApp configapp.Service, lf logfile.ILogFile, dom domainapp.DomainService) (IHTTPServer, func(), error) {
	cleanup := func() {
		log.Println("API cleanup")
	}

	srv, err := NewHTTPServer(new(configApp.Get().Hosts.API), lf.Get(), dom)
	if err != nil {
		return nil, cleanup, ge.Pin(err)
	}

	return srv, cleanup, nil
}
