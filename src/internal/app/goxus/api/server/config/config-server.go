package configserver

// HTTPServerConfig holds the HTTP server binding configuration.
type HTTPServerConfig struct {
	Host string `yaml:"host,omitempty"`
	Port string `yaml:"post,omitempty"`
}
