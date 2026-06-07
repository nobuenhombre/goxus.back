package settingsdomain

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/nobuenhombre/suikat/pkg/ge"

	"goxus/src/internal/pkg/db/goxus"
)

var (
	ErrSettingNotFound      = errors.New("setting not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrSettingsAccessDenied = errors.New("access denied")
)

// impl is the concrete implementation of Service with pure business logic.
type impl struct {
	repo *goxus.DbGoxusRepo
}

// New creates a new settings domain service.
func New(dbRepo *goxus.DbGoxusRepo) Service {
	return &impl{
		repo: dbRepo,
	}
}

// GetDefinitions returns all setting definitions enriched with type and group info.
func (s *impl) GetDefinitions() ([]*SettingsDefinition, error) {
	settings, err := s.repo.Setting.GetAll()
	if err != nil {
		return nil, ge.Pin(err)
	}

	result := make([]*SettingsDefinition, 0, len(settings))
	for _, setting := range settings {
		settingType, err := s.repo.SettingsType.GetSettingsTypeByID(setting.TypeID)
		if err != nil {
			return nil, ge.Pin(fmt.Errorf("resolve type %d for setting %d: %w", setting.TypeID, setting.ID, err))
		}

		settingGroup, err := s.repo.SettingsGroup.GetSettingsGroupByID(setting.GroupID)
		if err != nil {
			return nil, ge.Pin(fmt.Errorf("resolve group %d for setting %d: %w", setting.GroupID, setting.ID, err))
		}

		result = append(result, &SettingsDefinition{
			ID:   setting.ID,
			Type: settingType.Name,
			Group: func() string {
				if settingGroup.Name != "" {
					return settingGroup.Name
				}
				return ""
			}(),
			Name: setting.Name,
			Description: func() string {
				if setting.Description.Valid {
					return setting.Description.String
				}
				return ""
			}(),
			AvailableValues: setting.AvailableValues,
			DefaultValue:    setting.DefaultValue,
		})
	}

	return result, nil
}

// GetUserSettings returns all settings with user-specific values.
// Only returns settings where the user has a stored value.
func (s *impl) GetUserSettings(userID int64) ([]*UserSetting, error) {
	// Verify user exists
	_, err := s.repo.User.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ge.Pin(ErrUserNotFound)
		}
		return nil, ge.Pin(err)
	}

	allSettings, err := s.repo.Setting.GetAll()
	if err != nil {
		return nil, ge.Pin(err)
	}

	// Find user-specific settings by settings_id
	allUserSettings, err := s.repo.UsersSetting.GetAll()
	if err != nil {
		return nil, ge.Pin(err)
	}

	userSettingsMap := make(map[int64]*goxus.UsersSetting)
	for _, us := range allUserSettings {
		if us.UserID == userID {
			userSettingsMap[us.SettingsID] = us
		}
	}

	result := make([]*UserSetting, 0, len(allSettings))
	for _, setting := range allSettings {
		us, hasValue := userSettingsMap[setting.ID]
		if !hasValue {
			continue
		}

		settingType, err := s.repo.SettingsType.GetSettingsTypeByID(setting.TypeID)
		if err != nil {
			return nil, ge.Pin(fmt.Errorf("resolve type %d for setting %d: %w", setting.TypeID, setting.ID, err))
		}

		settingGroup, err := s.repo.SettingsGroup.GetSettingsGroupByID(setting.GroupID)
		if err != nil {
			return nil, ge.Pin(fmt.Errorf("resolve group %d for setting %d: %w", setting.GroupID, setting.ID, err))
		}

		result = append(result, &UserSetting{
			UserSettingsID: us.ID,
			SettingsID:     setting.ID,
			Type:           settingType.Name,
			Group: func() string {
				if settingGroup.Name != "" {
					return settingGroup.Name
				}
				return ""
			}(),
			Name: setting.Name,
			Description: func() string {
				if setting.Description.Valid {
					return setting.Description.String
				}
				return ""
			}(),
			AvailableValues: setting.AvailableValues,
			Value:           us.Value,
		})
	}

	return result, nil
}

// UpsertUserSetting creates or updates a user-specific setting value.
// Returns (created bool, error). created=true means a new row was inserted.
func (s *impl) UpsertUserSetting(userID, settingsID int64, value goxus.JSON) (bool, error) {
	// Verify user exists
	_, err := s.repo.User.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ge.Pin(ErrUserNotFound)
		}
		return false, ge.Pin(err)
	}

	// Verify setting exists
	_, err = s.repo.Setting.GetSettingByID(settingsID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ge.Pin(fmt.Errorf("settings_id '%d': %w", settingsID, ErrSettingNotFound))
		}
		return false, ge.Pin(err)
	}

	userSetting, err := s.repo.UsersSetting.GetUsersSettingByUserIDSettingsID(userID, settingsID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No existing setting — create new one below
			userSetting = nil
		} else {
			return false, ge.Pin(err)
		}
	}

	if userSetting != nil {
		userSetting.Value = value
		err = s.repo.UsersSetting.Save(userSetting)
		if err != nil {
			return false, ge.Pin(err)
		}
		return false, nil
	}

	// Create new
	us := &goxus.UsersSetting{
		UserID:     userID,
		SettingsID: settingsID,
		Value:      value,
	}
	err = s.repo.UsersSetting.Save(us)
	if err != nil {
		return false, ge.Pin(err)
	}

	return true, nil
}
