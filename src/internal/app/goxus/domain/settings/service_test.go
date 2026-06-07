package settingsdomain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goxus/src/internal/pkg/db/goxus"
)

// ---------------------------------------------------------------------------
// GetDefinitions
// ---------------------------------------------------------------------------

// TestGetDefinitions_Success verifies all setting definitions are returned
// with their type and group names resolved.
func TestGetDefinitions_Success(t *testing.T) {
	fx := setupTest(t)

	defs, err := fx.svc.GetDefinitions()
	require.NoError(t, err)
	require.NotEmpty(t, defs, "should have at least one seeded setting")

	first := defs[0]
	assert.NotZero(t, first.ID)
	assert.NotEmpty(t, first.Type, "type name should be resolved")
	assert.NotEmpty(t, first.Group, "group name should be resolved")
	assert.NotEmpty(t, first.Name, "setting name should be present")
	assert.NotNil(t, first.DefaultValue, "default_value should be present")
}

// TestGetDefinitions_ThemeSetting verifies the seeded Theme setting shape.
func TestGetDefinitions_ThemeSetting(t *testing.T) {
	fx := setupTest(t)

	defs, err := fx.svc.GetDefinitions()
	require.NoError(t, err)

	var theme *SettingsDefinition
	for _, d := range defs {
		if d.Name == "Theme" {
			theme = d
			break
		}
	}
	require.NotNil(t, theme, "Theme setting should be seeded")

	assert.Equal(t, "listRadios", theme.Type)
	assert.Equal(t, "Appearance", theme.Group)
	assert.Equal(t, "Select the theme for the dashboard.", theme.Description)
	assert.NotNil(t, theme.AvailableValues)
	assert.NotNil(t, theme.DefaultValue)
}

// ---------------------------------------------------------------------------
// GetUserSettings
// ---------------------------------------------------------------------------

// TestGetUserSettings_Success verifies a user with seeded settings gets them.
func TestGetUserSettings_Success(t *testing.T) {
	fx := setupTest(t)

	settings, err := fx.svc.GetUserSettings(fx.seedUserID)
	require.NoError(t, err)
	require.NotEmpty(t, settings, "seeded user should have at least one user setting")

	for _, s := range settings {
		assert.NotZero(t, s.UserSettingsID)
		assert.Equal(t, fx.seedSettingID, s.SettingsID)
		assert.Equal(t, fx.seedSettingType, s.Type)
		assert.Equal(t, fx.seedGroupName, s.Group)
		assert.NotEmpty(t, s.Name)
		assert.NotNil(t, s.Value, "value should be present for seeded user setting")
	}
}

// TestGetUserSettings_NotFound verifies ErrUserNotFound for non-existent user.
func TestGetUserSettings_NotFound(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.svc.GetUserSettings(99999)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// TestGetUserSettings_NoUserSpecificSetting returns empty when user has none.
func TestGetUserSettings_NoUserSpecificSetting(t *testing.T) {
	fx := setupTest(t)

	// Clean user settings for this test
	deleteAllUserSettings(t)
	t.Cleanup(func() { deleteAllUserSettings(t) })

	settings, err := fx.svc.GetUserSettings(fx.seedUserID)
	require.NoError(t, err)
	assert.Empty(t, settings)
}

// ---------------------------------------------------------------------------
// UpsertUserSetting
// ---------------------------------------------------------------------------

// TestUpsertUserSetting_Create verifies creating a new user setting.
func TestUpsertUserSetting_Create(t *testing.T) {
	fx := setupTest(t)

	// Remove the pre-seeded setting so we can test "create" path
	deleteAllUserSettings(t)
	t.Cleanup(func() { deleteAllUserSettings(t) })

	value := goxus.JSON{Data: map[string]any{"value": float64(2)}}

	created, err := fx.svc.UpsertUserSetting(fx.seedUserID, fx.seedSettingID, value)
	require.NoError(t, err)
	assert.True(t, created, "first upsert should create a new row")

	settings, err := fx.svc.GetUserSettings(fx.seedUserID)
	require.NoError(t, err)

	var found bool
	for _, s := range settings {
		if s.SettingsID == fx.seedSettingID {
			found = true
			assert.Equal(t, float64(2), s.Value.Data.(map[string]any)["value"])
		}
	}
	assert.True(t, found, "created setting should be findable")
}

// TestUpsertUserSetting_Update verifies updating an existing user setting.
func TestUpsertUserSetting_Update(t *testing.T) {
	fx := setupTest(t)

	// The seed setting exists — update it
	value := goxus.JSON{Data: map[string]any{"value": float64(2)}}

	created, err := fx.svc.UpsertUserSetting(fx.seedUserID, fx.seedSettingID, value)
	require.NoError(t, err)
	assert.False(t, created, "second upsert should update existing, not create")

	settings, err := fx.svc.GetUserSettings(fx.seedUserID)
	require.NoError(t, err)

	var found bool
	for _, s := range settings {
		if s.SettingsID == fx.seedSettingID {
			found = true
			assert.Equal(t, float64(2), s.Value.Data.(map[string]any)["value"])
		}
	}
	assert.True(t, found, "updated setting should be findable")
}

// TestUpsertUserSetting_SettingNotFound verifies error for invalid settings_id.
func TestUpsertUserSetting_SettingNotFound(t *testing.T) {
	fx := setupTest(t)

	value := goxus.JSON{Data: map[string]any{"value": float64(1)}}

	_, err := fx.svc.UpsertUserSetting(fx.seedUserID, 99999, value)
	require.Error(t, err)
	assert.ErrorContains(t, err, "setting not found")
}

// TestUpsertUserSetting_UserNotFound verifies error for invalid user_id.
func TestUpsertUserSetting_UserNotFound(t *testing.T) {
	fx := setupTest(t)

	value := goxus.JSON{Data: map[string]any{"value": float64(1)}}

	_, err := fx.svc.UpsertUserSetting(99999, fx.seedSettingID, value)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// TestUpsertUserSetting_VerifyPersistence verifies the upserted value survives re-read.
func TestUpsertUserSetting_VerifyPersistence(t *testing.T) {
	fx := setupTest(t)

	// Clean the seed, then create fresh
	deleteAllUserSettings(t)

	value := goxus.JSON{Data: map[string]any{"value": float64(3)}}
	created, err := fx.svc.UpsertUserSetting(fx.seedUserID, fx.seedSettingID, value)
	require.NoError(t, err)
	assert.True(t, created, "should create on clean slate")

	// Read back via service
	settings, err := fx.svc.GetUserSettings(fx.seedUserID)
	require.NoError(t, err)
	require.Len(t, settings, 1)
	assert.Equal(t, float64(3), settings[0].Value.Data.(map[string]any)["value"])
}
