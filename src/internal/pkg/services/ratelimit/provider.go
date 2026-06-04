package ratelimit

import (
	"log"
	"time"

	configapp "goxus/src/internal/app/goxus/config"

	"github.com/google/wire"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

// ProviderSet exports Wire providers for the ratelimit package.
var ProviderSet = wire.NewSet(
	ProvideRateLimiter,
)

// ProvideRateLimiter creates the rate limiter service from app config.
func ProvideRateLimiter(appConfig configapp.Service) (Service, func(), error) {
	cleanup := func() {
		log.Println("Rate limiter service cleanup")
	}

	cfg := appConfig.Get().RateLimit
	cfg.SetDefaults()

	if !cfg.Enabled {
		// disabled — provide a no-op rate limiter that always allows
		return New(Config{
			MaxAttempts: 0,
			Window:      0,
		}), cleanup, nil
	}

	window, err := time.ParseDuration(cfg.Window)
	if err != nil {
		return nil, cleanup, ge.Pin(err)
	}

	svc := New(Config{
		MaxAttempts: cfg.MaxAttempts,
		Window:      window,
	})

	return svc, cleanup, nil
}
