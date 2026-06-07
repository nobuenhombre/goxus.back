package settingsdomain

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"

	"goxus/src/internal/pkg/db/goxus"
	testpostgres "goxus/src/pkg/tests/postgres"
)

var globalRepo *goxus.DbGoxusRepo

type testFixtures struct {
	svc  Service
	repo *goxus.DbGoxusRepo

	// Seed helpers
	seedUserID      int64
	seedSettingID   int64
	seedSettingType string
	seedGroupName   string
}

func noopLog(_ string, _ time.Duration, _ ...any) {}

func TestMain(m *testing.M) {
	ctx := context.Background()

	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "../../../../../scripts/xo/goxus/migrations")

	container, dsn, err := testpostgres.StartPostgresContainer(ctx, testpostgres.VersionLatest)
	if err != nil {
		os.Exit(1)
	}

	migrateCmd := exec.Command("migrate", "-path", migrationsPath, "-database", dsn, "up")
	out, err := migrateCmd.CombinedOutput()
	if err != nil {
		os.Stderr.WriteString("migrate failed: " + err.Error() + "\n")
		os.Stderr.Write(out)
		container.Terminate(ctx)
		os.Exit(1)
	}

	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		os.Exit(1)
	}
	p, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		container.Terminate(ctx)
		os.Exit(1)
	}

	config := &pgxdb.Config{
		Host:     host,
		Port:     p.Port(),
		Name:     "testdb",
		User:     "testuser",
		Password: "testpass",
		SSLMode:  "disable",
	}

	dbRepo, err := goxus.NewDbGoxusRepository(config, noopLog)
	if err != nil {
		container.Terminate(ctx)
		os.Exit(1)
	}
	globalRepo = dbRepo

	code := m.Run()

	globalRepo.Close()
	container.Terminate(ctx)

	os.Exit(code)
}

// getFirstUser returns the first seeded user ID.
func getFirstUser(t *testing.T) int64 {
	users, err := globalRepo.User.GetAll()
	require.NoError(t, err)
	if len(users) == 0 {
		// Re-seed a user if migrations were clean somewhere
		now := time.Now()
		u := &goxus.User{
			Name:      "TestUser",
			Email:     "testuser@seed.com",
			Password:  "seedpass",
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := globalRepo.User.Save(u)
		require.NoError(t, err)
		return u.ID
	}
	return users[0].ID
}

// getFirstSettingID returns the first seeded setting definition ID.
func getFirstSettingID(t *testing.T) int64 {
	settings, err := globalRepo.Setting.GetAll()
	require.NoError(t, err)
	if len(settings) == 0 {
		t.Fatal("no settings found — migrations did not seed the settings table")
	}
	return settings[0].ID
}

// getFirstSettingType returns the type name of the first setting definition.
func getFirstSettingType(t *testing.T) string {
	setting, err := globalRepo.Setting.GetSettingByID(getFirstSettingID(t))
	require.NoError(t, err)
	st, err := globalRepo.SettingsType.GetSettingsTypeByID(setting.TypeID)
	require.NoError(t, err)
	return st.Name
}

// getFirstGroupName returns the group name of the first setting.
func getFirstGroupName(t *testing.T) string {
	setting, err := globalRepo.Setting.GetSettingByID(getFirstSettingID(t))
	require.NoError(t, err)
	sg, err := globalRepo.SettingsGroup.GetSettingsGroupByID(setting.GroupID)
	require.NoError(t, err)
	return sg.Name
}

// deleteAllUserSettings removes all rows from users_settings for a clean slate.
func deleteAllUserSettings(t *testing.T) {
	all, err := globalRepo.UsersSetting.GetAll()
	require.NoError(t, err)
	for _, us := range all {
		err := globalRepo.UsersSetting.Delete(us)
		require.NoError(t, err)
	}
}

// createUserSetting creates a users_settings row directly via the repo.
func createUserSetting(t *testing.T, userID, settingsID int64, value goxus.JSON) {
	us := &goxus.UsersSetting{
		UserID:     userID,
		SettingsID: settingsID,
		Value:      value,
	}
	err := globalRepo.UsersSetting.Save(us)
	require.NoError(t, err)
}

// setupTest returns isolated fixtures with fresh seed data for each test.
// Each call ensures:
//   - users_settings is empty (no cross-test pollution)
//   - seed user, settings_type, settings_group, and settings exist
//   - ONE standard seed users_setting is pre-populated
func setupTest(t *testing.T) testFixtures {
	t.Helper()

	svc := New(globalRepo)

	// Clean all user settings first
	deleteAllUserSettings(t)

	// Now create a fresh seed user setting for the first user
	userID := getFirstUser(t)
	settingID := getFirstSettingID(t)

	// Seed value matching the migration: {"value": 1}
	seedValue := goxus.JSON{Data: map[string]any{"value": float64(1)}}
	createUserSetting(t, userID, settingID, seedValue)

	t.Cleanup(func() {
		// Clean up any user settings the test may have created
		deleteAllUserSettings(t)
	})

	return testFixtures{
		svc:             svc,
		repo:            globalRepo,
		seedUserID:      userID,
		seedSettingID:   settingID,
		seedSettingType: getFirstSettingType(t),
		seedGroupName:   getFirstGroupName(t),
	}
}
