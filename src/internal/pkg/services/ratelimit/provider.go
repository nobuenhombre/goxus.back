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

	rlCfg := appConfig.Get().RateLimit

	if !rlCfg.Enabled {
		// disabled — provide a no-op rate limiter that always allows
		return New(Config{
			MaxAttempts: 0,
			Window:      0,
		}), cleanup, nil
	}

	window, err := time.ParseDuration(rlCfg.Window)
	if err != nil {
		return nil, cleanup, ge.Pin(err)
	}

	if rlCfg.MaxAttempts <= 0 {
		rlCfg.MaxAttempts = 5
	}
	if window <= 0 {
		window = 5 * time.Minute
	}

	svc := New(Config{
		MaxAttempts: rlCfg.MaxAttempts,
		Window:      window,
	})

	return svc, cleanup, nil
}
