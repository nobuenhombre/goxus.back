package userdomain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goxus/src/internal/pkg/services/rbac/permission"
	"goxus/src/internal/pkg/services/rbac/role"
)

// actorCtx returns a context with the given actor ID set.
func actorCtx(actorID int64) context.Context {
	return WithActorID(context.Background(), actorID)
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func TestAuthorizedCreate_NoActorID(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.svc.Create(context.Background(), "Alice", "a@b.com", "pass")
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedCreate_NoPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	// Permission exists in DB but actor has no role with it
	err := fx.rbacSvc.CreatePermission(permission.UserAdd, permission.UserAdd)
	require.NoError(t, err)

	_, err = fx.svc.Create(actorCtx(actorID), "Alice", "a@b.com", "pass")
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedCreate_WithPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserAdd)

	user, err := fx.svc.Create(actorCtx(actorID), "Alice", "a@b.com", "pass")
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

func TestAuthorizedList_NoActorID(t *testing.T) {
	fx := setupTest(t)

	_, _, err := fx.svc.List(context.Background(), 0, 0)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedList_NoPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	err := fx.rbacSvc.CreatePermission(permission.UserView, permission.UserView)
	require.NoError(t, err)

	_, _, err = fx.svc.List(actorCtx(actorID), 0, 0)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedList_WithPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserView)

	users, _, err := fx.svc.List(actorCtx(actorID), 0, 0)
	require.NoError(t, err)
	// List returns the actor user (created by createActor)
	assert.Len(t, users, 1)
}

// ---------------------------------------------------------------------------
// GetByID
// ---------------------------------------------------------------------------

func TestAuthorizedGetByID_NoActorID(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.svc.GetByID(context.Background(), 1)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedGetByID_NoPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	err := fx.rbacSvc.CreatePermission(permission.UserView, permission.UserView)
	require.NoError(t, err)

	// actorID != id (99999 vs actorID), no permission
	_, err = fx.svc.GetByID(actorCtx(actorID), 99999)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedGetByID_SelfRead(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	// actorID == id — should bypass permission check
	got, err := fx.svc.GetByID(actorCtx(actorID), actorID)
	require.NoError(t, err)
	assert.Equal(t, actorID, got.ID)
}

func TestAuthorizedGetByID_WithPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserView)

	// Different user, but has user_view permission
	created, err := fx.raw.Create(context.Background(), "Bob", "bob@test.com", "pass")
	require.NoError(t, err)

	got, err := fx.svc.GetByID(actorCtx(actorID), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestAuthorizedUpdate_NoActorID(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.svc.Update(context.Background(), 1, "X", "x@y.com")
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedUpdate_NoPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	err := fx.rbacSvc.CreatePermission(permission.UserEdit, permission.UserEdit)
	require.NoError(t, err)

	_, err = fx.svc.Update(actorCtx(actorID), actorID, "X", "x@y.com")
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedUpdate_WithPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserEdit)

	updated, err := fx.svc.Update(actorCtx(actorID), actorID, "Updated", "actor-new@test.com")
	require.NoError(t, err)
	assert.Equal(t, "Updated", updated.Name)
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestAuthorizedDelete_NoActorID(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.Delete(context.Background(), 1)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedDelete_NoPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	err := fx.rbacSvc.CreatePermission(permission.UserDelete, permission.UserDelete)
	require.NoError(t, err)

	err = fx.svc.Delete(actorCtx(actorID), 99999)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedDelete_WithPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserDelete)

	// Create a disposable user via raw (no auth) then delete via svc
	created, err := fx.raw.Create(context.Background(), "DeleteMe", "delete@test.com", "pass")
	require.NoError(t, err)

	err = fx.svc.Delete(actorCtx(actorID), created.ID)
	require.NoError(t, err)
}

func TestAuthorizedDelete_SelfDelete(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserDelete)

	// User tries to delete themselves — should fail
	err := fx.svc.Delete(actorCtx(actorID), actorID)
	require.Error(t, err)
	assert.ErrorContains(t, err, "cannot delete yourself")
}

// ---------------------------------------------------------------------------
// Restore
// ---------------------------------------------------------------------------

func TestAuthorizedRestore_NoActorID(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.svc.Restore(context.Background(), 1)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedRestore_NoPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	err := fx.rbacSvc.CreatePermission(permission.UserDelete, permission.UserDelete)
	require.NoError(t, err)

	_, err = fx.svc.Restore(actorCtx(actorID), 99999)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedRestore_WithPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserDelete)

	// Create a user, soft-delete via raw, then restore via svc
	created, err := fx.raw.Create(context.Background(), "RestoreMe", "restore@test.com", "pass")
	require.NoError(t, err)
	err = fx.raw.Delete(context.Background(), created.ID)
	require.NoError(t, err)

	restored, err := fx.svc.Restore(actorCtx(actorID), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, restored.ID)
	assert.False(t, restored.DeletedAt.Valid)
}

// ---------------------------------------------------------------------------
// Login — no permission check (public)
// ---------------------------------------------------------------------------

func TestAuthorizedLogin_NoActorID(t *testing.T) {
	fx := setupTest(t)

	// Login is public — should pass through even without actor
	_, err := fx.raw.Create(context.Background(), "LoginUser", "login@test.com", "pass")
	require.NoError(t, err)

	user, token, err := fx.svc.Login(context.Background(), "login@test.com", "pass")
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, token)
}

func TestAuthorizedLogin_WrongPassword(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "LoginUser", "login@test.com", "pass")
	require.NoError(t, err)

	_, _, err = fx.svc.Login(context.Background(), "login@test.com", "wrong")
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

// ---------------------------------------------------------------------------
// Logout — no permission check (token identifies itself)
// ---------------------------------------------------------------------------

func TestAuthorizedLogout_NoActorID(t *testing.T) {
	fx := setupTest(t)

	// Logout is public — should pass through even without actor
	_, err := fx.raw.Create(context.Background(), "LogoutUser", "logout@test.com", "pass")
	require.NoError(t, err)

	_, token, err := fx.raw.Login(context.Background(), "logout@test.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	err = fx.svc.Logout(context.Background(), token.Token)
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// ValidateToken — no permission check (public entry point)
// ---------------------------------------------------------------------------

func TestAuthorizedValidateToken_NoActorID(t *testing.T) {
	fx := setupTest(t)

	// ValidateToken is public — should pass through even without actor
	_, err := fx.raw.Create(context.Background(), "ValUser", "val@test.com", "pass")
	require.NoError(t, err)

	_, token, err := fx.raw.Login(context.Background(), "val@test.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	user, tok, err := fx.svc.ValidateToken(context.Background(), token.Token)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, tok)
	assert.Equal(t, token.Token, tok.Token)
}

func TestAuthorizedValidateToken_TokenNotFound(t *testing.T) {
	fx := setupTest(t)

	_, _, err := fx.svc.ValidateToken(context.Background(), "no-such-token")
	require.Error(t, err)
	assert.ErrorContains(t, err, "token not found")
}

// ---------------------------------------------------------------------------
// GetRoles
// ---------------------------------------------------------------------------

func TestAuthorizedGetRoles_NoActorID(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.svc.GetRoles(context.Background(), 1)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedGetRoles_NoPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	err := fx.rbacSvc.CreatePermission(permission.UserRoleView, permission.UserRoleView)
	require.NoError(t, err)

	_, err = fx.svc.GetRoles(actorCtx(actorID), 1)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedGetRoles_WithPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserRoleView)

	roles, err := fx.svc.GetRoles(actorCtx(actorID), actorID)
	require.NoError(t, err)
	// Actor was granted the "testrole" role — returns 1 role
	require.Len(t, roles, 1)
	assert.Equal(t, "testrole", roles[0].Slug)
}

// ---------------------------------------------------------------------------
// AssignRole
// ---------------------------------------------------------------------------

func TestAuthorizedAssignRole_NoActorID(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.AssignRole(context.Background(), 1, role.Admin)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedAssignRole_NoPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	err := fx.rbacSvc.CreatePermission(permission.UserRoleAdd, permission.UserRoleAdd)
	require.NoError(t, err)

	err = fx.svc.AssignRole(actorCtx(actorID), 1, role.Admin)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedAssignRole_WithPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserRoleAdd)

	// Create a role via rbacSvc directly first
	err := fx.rbacSvc.CreateRole("Admin", role.Admin)
	require.NoError(t, err)

	err = fx.svc.AssignRole(actorCtx(actorID), actorID, role.Admin)
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// RevokeRole
// ---------------------------------------------------------------------------

func TestAuthorizedRevokeRole_NoActorID(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.RevokeRole(context.Background(), 1, role.Admin)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedRevokeRole_NoPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	err := fx.rbacSvc.CreatePermission(permission.UserRoleDelete, permission.UserRoleDelete)
	require.NoError(t, err)

	err = fx.svc.RevokeRole(actorCtx(actorID), 1, role.Admin)
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedRevokeRole_WithPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserRoleDelete)

	// Assign a role first via raw (no auth)
	err := fx.rbacSvc.CreateRole("Temp", "temp")
	require.NoError(t, err)
	err = fx.raw.AssignRole(context.Background(), actorID, "temp")
	require.NoError(t, err)

	err = fx.svc.RevokeRole(actorCtx(actorID), actorID, "temp")
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// DeleteExpiredTokens — pass-through (no RBAC)
// ---------------------------------------------------------------------------

func TestAuthorizedDeleteExpiredTokens_Passthrough(t *testing.T) {
	fx := setupTest(t)

	// DeleteExpiredTokens is a system operation — no RBAC check, should pass through
	err := fx.svc.DeleteExpiredTokens(context.Background(), 30)
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// UpdatePassword — requires user_edit permission
// ---------------------------------------------------------------------------

func TestAuthorizedUpdatePassword_NoActorID(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.UpdatePassword(context.Background(), 1, "newpass")
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedUpdatePassword_NoPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)

	err := fx.rbacSvc.CreatePermission(permission.UserEdit, permission.UserEdit)
	require.NoError(t, err)

	err = fx.svc.UpdatePassword(actorCtx(actorID), actorID, "newpass")
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

func TestAuthorizedUpdatePassword_WithPermission(t *testing.T) {
	fx := setupTest(t)
	actorID := createActor(t, fx)
	grantPermission(t, fx, actorID, permission.UserEdit)

	err := fx.svc.UpdatePassword(actorCtx(actorID), actorID, "newpass")
	require.NoError(t, err)
}
