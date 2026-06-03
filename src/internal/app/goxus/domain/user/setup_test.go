package userdomain

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
	"github.com/stretchr/testify/require"

	"goxus/src/internal/pkg/db/goxus"
	"goxus/src/internal/pkg/services/rbac"
	testpostgres "goxus/src/pkg/tests/postgres"
)

var (
	globalRepo    *goxus.DbGoxusRepo
	globalRbacSvc rbac.Service
)

type testFixtures struct {
	svc     Service // top-level service with auth decorator
	raw     Service // raw impl without auth decorator
	repo    *goxus.DbGoxusRepo
	rbacSvc rbac.Service
}

func noopLog(_ string, _ time.Duration, _ ...any) {}

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Resolve migrations path relative to this file
	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "../../../../../scripts/xo/goxus/migrations")

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

	// 5. Create global RBAC service
	globalRbacSvc = rbac.New(globalRepo)

	// 6. Run all tests
	code := m.Run()

	// 7. Cleanup
	globalRepo.Close()
	container.Terminate(ctx)

	os.Exit(code)
}

// truncateAll cleans all test data between tests.
// Order respects FK constraints: child tables first.
func truncateAll(t *testing.T) {
	// 1. Delete all tokens
	allTokens, err := globalRepo.UsersToken.GetAll()
	require.NoError(t, err)
	for _, tok := range allTokens {
		err := globalRepo.UsersToken.Delete(tok)
		require.NoError(t, err)
	}

	// 2. Delete all users
	allUsers, err := globalRepo.User.GetAll()
	require.NoError(t, err)
	for _, u := range allUsers {
		err := globalRepo.User.Delete(u)
		require.NoError(t, err)
	}

	// 3. Delete RBAC join tables first (child tables)
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

	// 4. Delete RBAC parent tables
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

// setupTest returns isolated fixtures for each test.
// Each test runs against a clean slate (all tables truncated).
func setupTest(t *testing.T) testFixtures {
	truncateAll(t)

	rawSvc := New(globalRepo, globalRbacSvc)
	svc := NewAuthorized(rawSvc, globalRbacSvc)

	t.Cleanup(func() {
		truncateAll(t)
	})

	return testFixtures{
		svc:     svc,
		raw:     rawSvc,
		repo:    globalRepo,
		rbacSvc: globalRbacSvc,
	}
}

// createActor creates a new user in the DB and returns its ID for use as an actor.
func createActor(t *testing.T, fx testFixtures) int64 {
	now := time.Now()
	actor := &goxus.User{
		Name:      "ActorUser",
		Email:     fmt.Sprintf("actor-%d@test.com", now.UnixNano()),
		Password:  "actorpass",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := fx.repo.User.Save(actor)
	require.NoError(t, err)
	return actor.ID
}

// grantPermission creates the permission and a role, links them,
// and assigns the role to the given actor.
func grantPermission(t *testing.T, fx testFixtures, actorID int64, permSlug string) {
	err := fx.rbacSvc.CreatePermission(permSlug, permSlug)
	require.NoError(t, err)
	err = fx.rbacSvc.CreateRole("TestRole", "testrole")
	require.NoError(t, err)
	err = fx.rbacSvc.AssignPermissionsToRole("testrole", []string{permSlug})
	require.NoError(t, err)
	err = fx.rbacSvc.AssignRoleToUser(actorID, "testrole")
	require.NoError(t, err)
}
