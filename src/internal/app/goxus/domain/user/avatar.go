package userdomain

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	_ "golang.org/x/image/webp"

	"github.com/nobuenhombre/suikat/pkg/ge"
)

//go:embed default_avatar.webp
var embeddedDefaultAvatar []byte

const (
	// AvatarWidth is the required avatar image width in pixels.
	AvatarWidth = 460
	// AvatarHeight is the required avatar image height in pixels.
	AvatarHeight = 460
	// DefaultAvatarFilename is the filename of the default avatar.
	DefaultAvatarFilename = "user-avatar-default.webp"
	// avatarContentType is the content type returned for avatar images.
	avatarContentType = "image/webp"
)

// avatarFilename returns the avatar filename for a given user ID.
func avatarFilename(userID int64) string {
	return fmt.Sprintf("user-avatar-%d.webp", userID)
}

// avatarPath returns the full avatar file path for a given user ID.
func avatarPath(avatarsDir string, userID int64) string {
	return filepath.Join(avatarsDir, avatarFilename(userID))
}

// defaultAvatarPath returns the full path to the default avatar file.
func defaultAvatarPath(avatarsDir string) string {
	return filepath.Join(avatarsDir, DefaultAvatarFilename)
}

// validateImageDimensions checks that the image data decodes to exactly 460x460 pixels.
func validateImageDimensions(data []byte) error {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return ge.Pin(fmt.Errorf("%w: %v", ErrInvalidImageFormat, err))
	}

	if cfg.Width != AvatarWidth || cfg.Height != AvatarHeight {
		return ge.Pin(fmt.Errorf("%w: got %dx%d", ErrInvalidImageSize, cfg.Width, cfg.Height))
	}

	return nil
}

// GetAvatar returns the avatar image for a user.
// Returns the user's custom avatar if it exists, otherwise the default avatar.
// If the default avatar file does not exist on disk, the embedded default is used.
func (s *impl) GetAvatar(_ context.Context, userID int64) ([]byte, string, error) {
	// Try user-specific avatar first
	path := avatarPath(s.cfg.Storage.AvatarsDir, userID)
	data, err := os.ReadFile(path)
	if err == nil {
		return data, avatarContentType, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, "", ge.Pin(fmt.Errorf("reading avatar for user %d: %w", userID, err))
	}

	// Fall back to default avatar on disk
	defaultPath := defaultAvatarPath(s.cfg.Storage.AvatarsDir)
	data, err = os.ReadFile(defaultPath)
	if err == nil {
		return data, avatarContentType, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, "", ge.Pin(fmt.Errorf("reading default avatar: %w", err))
	}

	// Default avatar not on disk — write the embedded default for next time
	if err := os.MkdirAll(s.cfg.Storage.AvatarsDir, 0o755); err != nil {
		return nil, "", ge.Pin(fmt.Errorf("creating avatars directory: %w", err))
	}
	if err := os.WriteFile(defaultPath, embeddedDefaultAvatar, 0o644); err != nil {
		return nil, "", ge.Pin(fmt.Errorf("writing default avatar: %w", err))
	}

	return embeddedDefaultAvatar, avatarContentType, nil
}

// UploadAvatar saves an avatar image for a user.
// Validates that the image is exactly 460x460 pixels before saving.
func (s *impl) UploadAvatar(_ context.Context, userID int64, data []byte) error {
	if len(data) == 0 {
		return ge.Pin(fmt.Errorf("empty image data: %w", ErrInvalidImageFormat))
	}

	// Validate image dimensions
	err := validateImageDimensions(data)
	if err != nil {
		return ge.Pin(err)
	}

	// Ensure the avatars directory exists
	err = os.MkdirAll(s.cfg.Storage.AvatarsDir, 0o755)
	if err != nil {
		return ge.Pin(fmt.Errorf("creating avatars directory: %w", err))
	}

	// Write the avatar file
	path := avatarPath(s.cfg.Storage.AvatarsDir, userID)
	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		return ge.Pin(fmt.Errorf("writing avatar for user %d: %w", userID, err))
	}

	return nil
}

// DeleteAvatar removes the custom avatar for a user.
// After deletion, GetAvatar will return the default avatar.
func (s *impl) DeleteAvatar(_ context.Context, userID int64) error {
	path := avatarPath(s.cfg.Storage.AvatarsDir, userID)

	err := os.Remove(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// No custom avatar to delete — not an error
			return nil
		}
		return ge.Pin(fmt.Errorf("deleting avatar for user %d: %w", userID, err))
	}

	return nil
}
