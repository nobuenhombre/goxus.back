// Package configstorage provides configuration for file storage paths.
// This is a cross-cutting config that can be extended with new path fields
// as the application grows (other images, uploaded files, etc.).
package configstorage

// DefaultAvatarsDir is the default directory for user avatar images.
const DefaultAvatarsDir = "data/production/img/users/avatars"

// StorageConfig holds file system path configuration.
// All paths are relative to the application working directory.
type StorageConfig struct {
	AvatarsDir string `yaml:"avatars_dir,omitempty"`
}

// SetDefaults applies default values for zero-value fields.
func (c *StorageConfig) SetDefaults() {
	if c.AvatarsDir == "" {
		c.AvatarsDir = DefaultAvatarsDir
	}
}
