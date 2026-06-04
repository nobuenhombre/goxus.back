package configratelimit

const (
	DefaultMaxAttempts int = 5
	DefaultWindow          = "5m"
)

// LoginRateLimitConfig holds rate-limiting configuration for the login endpoint.
type LoginRateLimitConfig struct {
	Enabled     bool   `yaml:"enabled"`
	MaxAttempts int    `yaml:"max_attempts"`
	Window      string `yaml:"window"` // duration string, e.g. "5m"
}

// SetDefaults apply defaults for zero-valued config
func (c *LoginRateLimitConfig) SetDefaults() {
	if c.MaxAttempts == 0 {
		c.MaxAttempts = DefaultMaxAttempts
	}

	if c.Window == "" {
		c.Window = DefaultWindow
	}
}
