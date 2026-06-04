package configcleantokens

const DefaultTokenTTLDays = 7

type TokenCleanupConfig struct {
	TTLDays int `yaml:"ttl_days,omitempty"` // days until token expiry from last_used_at
}

// SetDefaults apply defaults for zero-valued config
func (c *TokenCleanupConfig) SetDefaults() {
	if c.TTLDays == 0 {
		c.TTLDays = DefaultTokenTTLDays
	}
}
