package rbac

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
	"github.com/stretchr/testify/require"

	"goxus/src/internal/pkg/db/goxus"
	testpostgres "goxus/src/pkg/tests/postgres"
)

var (
	globalRepo     *goxus.DbGoxusRepo
	testUserID     int64
	migrationsPath string
)

type testFixtures struct {
	svc  Service
	repo *goxus.DbGoxusRepo
}

func noopLog(_ string, _ time.Duration, _ ...any) {}

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Resolve migrations path relative to this file
	_, filename, _, _ := runtime.Caller(0)
	migrationsPath = filepath.Join(filepath.Dir(filename), "../../../../scripts/xo/goxus/migrations")

	// 1. Start one Postgres container for all tests
	container, dsn, err := testpostgres.StartPostgresContainer(ctx, testpostgres.VersionLatest)
	if err != nil {
		os.Exit(1)
	}

	// 2. Get connection params
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

	// 3. Apply migrations
	migrateCmd := exec.Command(
		"migrate",
		"-path", migrationsPath,
		"-database", dsn,
		"up",
	)
	out, err := migrateCmd.CombinedOutput()
	if err != nil {
		os.Stderr.WriteString("migrate failed: " + err.Error() + "\n")
		os.Stderr.Write(out)
		os.Stderr.WriteString("\n")
		container.Terminate(ctx)
		os.Exit(1)
	}

	// 4. Create DbGoxusRepo for tests
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

	// 5. Create a test user instead of relying on seeded ID=1
	testUser := &goxus.User{
		Name:     "TestUser",
		Email:    "testuser@example.com",
		Password: "testpass",
	}
	err = globalRepo.User.Save(testUser)
	if err != nil {
		globalRepo.Close()
		container.Terminate(ctx)
		os.Exit(1)
	}
	testUserID = testUser.ID

	// 6. Run all tests
	code := m.Run()

	// 7. Cleanup
	globalRepo.Close()
	container.Terminate(ctx)

	os.Exit(code)
}

// testUserID is set in TestMain from a dynamically created test user.
// seeded via migration 000002 is no longer relied upon.

// truncateRBAC cleans all RBAC tables between tests.
// Order respects FK constraints: child tables first.
func truncateRBAC(t *testing.T) {
	// 1. Child tables first: rbac_role_permissions, rbac_user_roles
	allRolePerms, err := globalRepo.RbacRolePermission.GetAll()
	require.NoError(t, err)
	for _, rp := range allRolePerms {
		err := globalRepo.RbacRolePermission.Delete(rp)
		require.NoError(t, err)
	}

	allUserRoles, err := globalRepo.RbacUserRole.GetAll()
	require.NoError(t, err)
	for _, ur := range allUserRoles {
		err := globalRepo.RbacUserRole.Delete(ur)
		require.NoError(t, err)
	}

	// 2. Parent tables: rbac_permissions, rbac_roles
	allPerms, err := globalRepo.RbacPermission.GetAll()
	require.NoError(t, err)
	for _, p := range allPerms {
		err := globalRepo.RbacPermission.Delete(p)
		require.NoError(t, err)
	}

	allRoles, err := globalRepo.RbacRole.GetAll()
	require.NoError(t, err)
	for _, r := range allRoles {
		err := globalRepo.RbacRole.Delete(r)
		require.NoError(t, err)
	}
}

func setupTest(t *testing.T) testFixtures {
	truncateRBAC(t)

	svc := New(globalRepo)

	t.Cleanup(func() {
		truncateRBAC(t)
	})

	return testFixtures{
		svc:  svc,
		repo: globalRepo,
	}
}
