// Package settingsdomain provides the domain service for user settings management.
// It handles setting definitions and user-specific setting values.
package settingsdomain

import (
	"goxus/src/internal/pkg/db/goxus"
)

// SettingsDefinition represents a setting definition enriched with type and group info.
// This is the response for GET /api/v1/entity/settings.
type SettingsDefinition struct {
	ID              int64      `json:"id"`
	Type            string     `json:"type"`
	Group           string     `json:"group"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	AvailableValues goxus.JSON `json:"available_values"`
	DefaultValue    goxus.JSON `json:"default_value"`
}

// UserSetting represents a user's setting value enriched with setting definition.
// This is the response for GET /api/v1/entity/user/:id/settings.
type UserSetting struct {
	UserSettingsID  int64      `json:"user_settings_id"`
	SettingsID      int64      `json:"settings_id"`
	Type            string     `json:"type"`
	Group           string     `json:"group"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	AvailableValues goxus.JSON `json:"available_values"`
	Value           goxus.JSON `json:"value"`
}

// Service defines the settings domain operations.
type Service interface {
	// GetDefinitions returns all setting definitions enriched with type and group info.
	GetDefinitions() ([]*SettingsDefinition, error)

	// GetUserSettings returns all settings with user-specific values.
	// If a setting has no user-specific value, it is NOT included in the result
	// (frontend uses default_value from the definition).
	GetUserSettings(userID int64) ([]*UserSetting, error)

	// UpsertUserSetting creates or updates a user-specific setting value.
	// Returns (created bool, error). created=true means a new row was inserted.
	UpsertUserSetting(userID, settingsID int64, value goxus.JSON) (bool, error)
}
