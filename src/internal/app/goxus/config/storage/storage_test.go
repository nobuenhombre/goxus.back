package configstorage

import "testing"

func TestSetDefaults(t *testing.T) {
	t.Run("empty_gets_default", func(t *testing.T) {
		c := &StorageConfig{}
		c.SetDefaults()

		if c.AvatarsDir != DefaultAvatarsDir {
			t.Errorf(
				"SetDefaults() with empty AvatarsDir:\n expected %q,\n got      %q",
				DefaultAvatarsDir, c.AvatarsDir,
			)
		}
	})

	t.Run("non_empty_stays", func(t *testing.T) {
		customPath := "custom/img/path"
		c := &StorageConfig{AvatarsDir: customPath}
		c.SetDefaults()

		if c.AvatarsDir != customPath {
			t.Errorf(
				"SetDefaults() with non-empty AvatarsDir:\n expected %q,\n got      %q",
				customPath, c.AvatarsDir,
			)
		}
	})
}
