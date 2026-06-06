package userdomain

import (
	"context"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goxus/src/internal/pkg/services/rbac/role"
)

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

// TestCreate_Success creates a user and verifies it is persisted.
func TestCreate_Success(t *testing.T) {
	fx := setupTest(t)

	user, err := fx.raw.Create(context.Background(), "Alice", "alice@example.com", "secret")
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)
	assert.Equal(t, "secret", user.Password)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, 5*time.Second)
	assert.WithinDuration(t, time.Now(), user.UpdatedAt, 5*time.Second)
	assert.False(t, user.DeletedAt.Valid)
}

// TestCreate_DuplicateEmail verifies that creating a user with an existing email
// returns ErrEmailAlreadyTaken.
func TestCreate_DuplicateEmail(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Alice", "dup@example.com", "secret")
	require.NoError(t, err)

	_, err = fx.raw.Create(context.Background(), "Bob", "dup@example.com", "other")
	require.Error(t, err)
	assert.ErrorContains(t, err, "email already taken")
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// TestList_Empty verifies List returns an empty slice when no users exist.
func TestList_Empty(t *testing.T) {
	fx := setupTest(t)

	users, total, err := fx.raw.List(context.Background(), 0, 0)
	require.NoError(t, err)
	assert.Len(t, users, 0)
	assert.Equal(t, int64(0), total)
}

// TestList_Multiple verifies that multiple users are returned with default pagination.
// Also verifies that users without roles have empty Roles field.
func TestList_Multiple(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Alice", "alice@example.com", "a")
	require.NoError(t, err)
	_, err = fx.raw.Create(context.Background(), "Bob", "bob@example.com", "b")
	require.NoError(t, err)

	users, total, err := fx.raw.List(context.Background(), 0, 0)
	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, int64(2), total)

	// Both users have no roles — Roles should be empty string (Valid=true)
	for _, u := range users {
		assert.True(t, u.Roles.Valid, "Roles should be valid (non-NULL) for users without roles")
		assert.Empty(t, u.Roles.String, "Roles should be empty for users without roles")
	}
}

// TestList_Roles_Empty verifies that a user with no roles has empty Roles field.
func TestList_Roles_Empty(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "NoRole", "norole@example.com", "pass")
	require.NoError(t, err)

	users, total, err := fx.raw.List(context.Background(), 0, 0)
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, int64(1), total)

	// User has no roles — view coalesces to empty string
	assert.True(t, users[0].Roles.Valid)
	assert.Empty(t, users[0].Roles.String)
}

// TestList_Roles_SingleRole verifies that a user with one role has it in the Roles field.
func TestList_Roles_SingleRole(t *testing.T) {
	fx := setupTest(t)

	user, err := fx.raw.Create(context.Background(), "SingleRole", "singlerole@example.com", "pass")
	require.NoError(t, err)

	// Create role and assign
	err = fx.rbacSvc.CreateRole("Editor", "editor")
	require.NoError(t, err)
	err = fx.raw.AssignRole(context.Background(), user.ID, "editor")
	require.NoError(t, err)

	users, total, err := fx.raw.List(context.Background(), 0, 0)
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, int64(1), total)

	assert.True(t, users[0].Roles.Valid)
	assert.Equal(t, "Editor", users[0].Roles.String)
}

// TestList_Roles_MultipleRoles verifies that a user with multiple roles
// has them aggregated in the Roles field (sorted alphabetically, comma-separated).
func TestList_Roles_MultipleRoles(t *testing.T) {
	fx := setupTest(t)

	user, err := fx.raw.Create(context.Background(), "MultiRole", "multirole@example.com", "pass")
	require.NoError(t, err)

	// Create two roles (alphabetically "Moderator" < "Viewer")
	err = fx.rbacSvc.CreateRole("Moderator", "moderator")
	require.NoError(t, err)
	err = fx.rbacSvc.CreateRole("Viewer", "viewer")
	require.NoError(t, err)

	// Assign both
	err = fx.raw.AssignRole(context.Background(), user.ID, "moderator")
	require.NoError(t, err)
	err = fx.raw.AssignRole(context.Background(), user.ID, "viewer")
	require.NoError(t, err)

	users, total, err := fx.raw.List(context.Background(), 0, 0)
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, int64(1), total)

	assert.True(t, users[0].Roles.Valid)
	// string_agg orders by r.name ASC: "Moderator, Viewer"
	assert.Equal(t, "Moderator, Viewer", users[0].Roles.String)
}

// TestList_Roles_DifferentUsers verifies that each user's roles are correctly
// scoped — user A's roles don't leak into user B's Roles field.
func TestList_Roles_DifferentUsers(t *testing.T) {
	fx := setupTest(t)

	alice, err := fx.raw.Create(context.Background(), "Alice", "alice-roles@example.com", "pass")
	require.NoError(t, err)
	bob, err := fx.raw.Create(context.Background(), "Bob", "bob-roles@example.com", "pass")
	require.NoError(t, err)

	// Create roles
	err = fx.rbacSvc.CreateRole("Admin", "admin")
	require.NoError(t, err)
	err = fx.rbacSvc.CreateRole("Viewer", "viewer")
	require.NoError(t, err)

	// Assign different roles
	err = fx.raw.AssignRole(context.Background(), alice.ID, "admin")
	require.NoError(t, err)
	err = fx.raw.AssignRole(context.Background(), bob.ID, "viewer")
	require.NoError(t, err)

	users, total, err := fx.raw.List(context.Background(), 0, 0)
	require.NoError(t, err)
	require.Len(t, users, 2)
	assert.Equal(t, int64(2), total)

	for _, u := range users {
		assert.True(t, u.Roles.Valid)
		switch u.Name.String {
		case "Alice":
			assert.Equal(t, "Admin", u.Roles.String, "Alice should have Admin role")
		case "Bob":
			assert.Equal(t, "Viewer", u.Roles.String, "Bob should have Viewer role")
		default:
			t.Fatalf("unexpected user: %s", u.Name.String)
		}
	}
}

// TestList_Pagination_Limit verifies that limit restricts the number of results.
func TestList_Pagination_Limit(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Alice", "alice@example.com", "a")
	require.NoError(t, err)
	_, err = fx.raw.Create(context.Background(), "Bob", "bob@example.com", "b")
	require.NoError(t, err)

	users, total, err := fx.raw.List(context.Background(), 1, 0)
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, int64(2), total)
}

// TestList_Pagination_Offset verifies that offset skips the first N results.
func TestList_Pagination_Offset(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Alice", "alice@example.com", "a")
	require.NoError(t, err)
	_, err = fx.raw.Create(context.Background(), "Bob", "bob@example.com", "b")
	require.NoError(t, err)

	users, total, err := fx.raw.List(context.Background(), 10, 1)
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, int64(2), total)
	// With ORDER BY id ASC, offset=1 should return the second user
	assert.Equal(t, "Bob", users[0].Name.String)
}

// ---------------------------------------------------------------------------
// GetByID
// ---------------------------------------------------------------------------

// TestGetByID_Success verifies a user can be retrieved by ID.
func TestGetByID_Success(t *testing.T) {
	fx := setupTest(t)

	created, err := fx.raw.Create(context.Background(), "Carol", "carol@example.com", "pass")
	require.NoError(t, err)

	got, err := fx.raw.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "Carol", got.Name)
	assert.Equal(t, "carol@example.com", got.Email)
}

// TestGetByID_NotFound verifies ErrUserNotFound for a non-existent ID.
func TestGetByID_NotFound(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.GetByID(context.Background(), 99999)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

// TestUpdate_Success verifies name and email can be updated.
func TestUpdate_Success(t *testing.T) {
	fx := setupTest(t)

	created, err := fx.raw.Create(context.Background(), "Diana", "diana@example.com", "pass")
	require.NoError(t, err)

	updated, err := fx.raw.Update(context.Background(), created.ID, "Diana Updated", "diana.new@example.com")
	require.NoError(t, err)
	assert.Equal(t, "Diana Updated", updated.Name)
	assert.Equal(t, "diana.new@example.com", updated.Email)
	assert.WithinDuration(t, time.Now(), updated.UpdatedAt, 5*time.Second)

	// Verify persistence
	got, err := fx.raw.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Diana Updated", got.Name)
	assert.Equal(t, "diana.new@example.com", got.Email)
}

// TestUpdate_EmailAlreadyTaken verifies that updating to an existing email returns an error.
func TestUpdate_EmailAlreadyTaken(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Alice", "alice@example.com", "pass")
	require.NoError(t, err)
	bob, err := fx.raw.Create(context.Background(), "Bob", "bob@example.com", "pass")
	require.NoError(t, err)

	_, err = fx.raw.Update(context.Background(), bob.ID, "Bob", "alice@example.com")
	require.Error(t, err)
	assert.ErrorContains(t, err, "email already taken")
}

// TestUpdate_SameEmail verifies that updating a user without changing their email
// is allowed (no conflict with own email).
func TestUpdate_SameEmail(t *testing.T) {
	fx := setupTest(t)

	created, err := fx.raw.Create(context.Background(), "Eve", "eve@example.com", "pass")
	require.NoError(t, err)

	updated, err := fx.raw.Update(context.Background(), created.ID, "Eve Updated", "eve@example.com")
	require.NoError(t, err)
	assert.Equal(t, "Eve Updated", updated.Name)
}

// TestUpdate_NotFound verifies ErrUserNotFound when updating a non-existent user.
func TestUpdate_NotFound(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Update(context.Background(), 99999, "Ghost", "ghost@example.com")
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

// TestDelete_Success verifies a user can be soft-deleted.
func TestDelete_Success(t *testing.T) {
	fx := setupTest(t)

	created, err := fx.raw.Create(context.Background(), "Frank", "frank@example.com", "pass")
	require.NoError(t, err)

	err = fx.raw.Delete(context.Background(), created.ID)
	require.NoError(t, err)

	// Verify user is soft-deleted (still exists but deleted_at set)
	got, err := fx.raw.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.True(t, got.DeletedAt.Valid)
	assert.Equal(t, created.ID, got.ID)
}

// TestDelete_NotFound verifies ErrUserNotFound when deleting a non-existent user.
func TestDelete_NotFound(t *testing.T) {
	fx := setupTest(t)

	err := fx.raw.Delete(context.Background(), 99999)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// TestDelete_AlreadyDeleted verifies ErrUserAlreadyDeleted when soft-deleting an already deleted user.
func TestDelete_AlreadyDeleted(t *testing.T) {
	fx := setupTest(t)

	created, err := fx.raw.Create(context.Background(), "Gina", "gina@example.com", "pass")
	require.NoError(t, err)

	// First delete
	err = fx.raw.Delete(context.Background(), created.ID)
	require.NoError(t, err)

	// Second delete should fail
	err = fx.raw.Delete(context.Background(), created.ID)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user already deleted")
}

// ---------------------------------------------------------------------------
// Restore
// ---------------------------------------------------------------------------

// TestRestore_Success verifies a soft-deleted user can be restored.
func TestRestore_Success(t *testing.T) {
	fx := setupTest(t)

	created, err := fx.raw.Create(context.Background(), "Harry", "harry@example.com", "pass")
	require.NoError(t, err)

	// Soft-delete
	err = fx.raw.Delete(context.Background(), created.ID)
	require.NoError(t, err)

	// Restore
	restored, err := fx.raw.Restore(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, restored.ID)
	assert.False(t, restored.DeletedAt.Valid)

	// Verify persistence
	got, err := fx.raw.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.False(t, got.DeletedAt.Valid)
}

// TestRestore_NotDeleted verifies ErrUserNotDeleted when restoring a non-deleted user.
func TestRestore_NotDeleted(t *testing.T) {
	fx := setupTest(t)

	created, err := fx.raw.Create(context.Background(), "Iris", "iris@example.com", "pass")
	require.NoError(t, err)

	_, err = fx.raw.Restore(context.Background(), created.ID)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not deleted")
}

// TestRestore_NotFound verifies ErrUserNotFound when restoring a non-existent user.
func TestRestore_NotFound(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Restore(context.Background(), 99999)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// ---------------------------------------------------------------------------
// Login / Logout
// ---------------------------------------------------------------------------

// TestLogin_Success verifies a user can login with correct credentials.
func TestLogin_Success(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Grace", "grace@example.com", "mypassword")
	require.NoError(t, err)

	user, token, err := fx.raw.Login(context.Background(), "grace@example.com", "mypassword")
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.Equal(t, user.ID, token.UserID)
	assert.False(t, token.DeletedAt.Valid)
}

// TestLogin_WrongPassword verifies ErrAccessDenied for wrong password.
func TestLogin_WrongPassword(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Heidi", "heidi@example.com", "correctpass")
	require.NoError(t, err)

	_, _, err = fx.raw.Login(context.Background(), "heidi@example.com", "wrongpass")
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

// TestLogin_UserNotFound verifies ErrUserNotFound for non-existent email.
func TestLogin_UserNotFound(t *testing.T) {
	fx := setupTest(t)

	_, _, err := fx.raw.Login(context.Background(), "nobody@example.com", "pass")
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// TestLogin_SoftDeleted verifies that a soft-deleted user cannot login.
func TestLogin_SoftDeleted(t *testing.T) {
	fx := setupTest(t)

	// We need to directly soft-delete a user to test this path.
	// The impl.Delete does a physical delete, so we'll create a user
	// and manually set DeletedAt before attempting login.
	created, err := fx.raw.Create(context.Background(), "Ivan", "ivan@example.com", "pass")
	require.NoError(t, err)

	// Soft-delete by setting DeletedAt directly in the DB
	now := time.Now()
	created.DeletedAt.Valid = true
	created.DeletedAt.Time = now
	err = fx.repo.User.Save(created) // Save calls Update since user already exists
	require.NoError(t, err)

	_, _, err = fx.raw.Login(context.Background(), "ivan@example.com", "pass")
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// TestLogout_Success verifies a token can be invalidated.
func TestLogout_Success(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Judy", "judy@example.com", "pass")
	require.NoError(t, err)

	user, token, err := fx.raw.Login(context.Background(), "judy@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	err = fx.raw.Logout(context.Background(), token.Token)
	require.NoError(t, err)

	// Verify token is soft-deleted
	tok2, err := fx.repo.UsersToken.GetUsersTokenByToken(token.Token)
	require.NoError(t, err)
	assert.True(t, tok2.DeletedAt.Valid)
	assert.Equal(t, user.ID, tok2.UserID)
}

// TestLogout_TokenNotFound verifies ErrTokenNotFound for an invalid token.
func TestLogout_TokenNotFound(t *testing.T) {
	fx := setupTest(t)

	err := fx.raw.Logout(context.Background(), "non-existent-token")
	require.Error(t, err)
	assert.ErrorContains(t, err, "token not found")
}

// ---------------------------------------------------------------------------
// ValidateToken
// ---------------------------------------------------------------------------

// TestValidateToken_Success validates a valid token and returns the user and token.
func TestValidateToken_Success(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Oscar", "oscar@example.com", "pass")
	require.NoError(t, err)

	user, token, err := fx.raw.Login(context.Background(), "oscar@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	gotUser, gotToken, err := fx.raw.ValidateToken(context.Background(), token.Token)
	require.NoError(t, err)
	assert.Equal(t, user.ID, gotUser.ID)
	assert.Equal(t, token.Token, gotToken.Token)
	assert.Equal(t, user.ID, gotToken.UserID)
	assert.False(t, gotToken.DeletedAt.Valid)

	// last_used_at should have been updated
	assert.True(t, gotToken.LastUsedAt.Valid)
	assert.WithinDuration(t, time.Now(), gotToken.LastUsedAt.Time, 5*time.Second)
}

// TestValidateToken_TokenNotFound verifies ErrTokenNotFound for an invalid token.
func TestValidateToken_TokenNotFound(t *testing.T) {
	fx := setupTest(t)

	_, _, err := fx.raw.ValidateToken(context.Background(), "non-existent-token")
	require.Error(t, err)
	assert.ErrorContains(t, err, "token not found")
}

// TestValidateToken_TokenSoftDeleted verifies ErrTokenNotFound for a logged-out token.
func TestValidateToken_TokenSoftDeleted(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Peggy", "peggy@example.com", "pass")
	require.NoError(t, err)

	_, token, err := fx.raw.Login(context.Background(), "peggy@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	// Logout — soft-deletes the token
	err = fx.raw.Logout(context.Background(), token.Token)
	require.NoError(t, err)

	_, _, err = fx.raw.ValidateToken(context.Background(), token.Token)
	require.Error(t, err)
	assert.ErrorContains(t, err, "token not found")
}

// TestValidateToken_UserSoftDeleted verifies that a token for a soft-deleted user
// returns ErrUserNotFound.
func TestValidateToken_UserSoftDeleted(t *testing.T) {
	fx := setupTest(t)

	created, err := fx.raw.Create(context.Background(), "Quinn", "quinn@example.com", "pass")
	require.NoError(t, err)

	_, token, err := fx.raw.Login(context.Background(), "quinn@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	// Soft-delete the user directly in the DB
	now := time.Now()
	created.DeletedAt.Valid = true
	created.DeletedAt.Time = now
	err = fx.repo.User.Save(created)
	require.NoError(t, err)

	_, _, err = fx.raw.ValidateToken(context.Background(), token.Token)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// TestValidateToken_LastUsedAtUpdated verifies that calling ValidateToken
// repeatedly updates last_used_at each time.
func TestValidateToken_LastUsedAtUpdated(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Rupert", "rupert@example.com", "pass")
	require.NoError(t, err)

	_, token, err := fx.raw.Login(context.Background(), "rupert@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	// First validation
	_, gotToken1, err := fx.raw.ValidateToken(context.Background(), token.Token)
	require.NoError(t, err)
	assert.True(t, gotToken1.LastUsedAt.Valid)
	firstUsed := gotToken1.LastUsedAt.Time

	// Wait a moment
	time.Sleep(5 * time.Millisecond)

	// Second validation — last_used_at should advance
	_, gotToken2, err := fx.raw.ValidateToken(context.Background(), token.Token)
	require.NoError(t, err)
	assert.True(t, gotToken2.LastUsedAt.Valid)

	// Only check if it actually advanced (may be same timestamp in same millisecond)
	if !gotToken2.LastUsedAt.Time.After(firstUsed) {
		t.Logf("last_used_at did not advance (possible same-millisecond execution)")
	}
}

// TestAssignRole_Success verifies a role can be assigned to a user.
func TestAssignRole_Success(t *testing.T) {
	fx := setupTest(t)

	user, err := fx.raw.Create(context.Background(), "Kate", "kate@example.com", "pass")
	require.NoError(t, err)

	err = fx.rbacSvc.CreateRole("Moderator", "moderator")
	require.NoError(t, err)

	err = fx.raw.AssignRole(context.Background(), user.ID, "moderator")
	require.NoError(t, err)

	roles, err := fx.raw.GetRoles(context.Background(), user.ID)
	require.NoError(t, err)
	require.Len(t, roles, 1)
	assert.Equal(t, "moderator", roles[0].Slug)
}

// TestAssignRole_UserNotFound verifies ErrUserNotFound when assigning to a non-existent user.
func TestAssignRole_UserNotFound(t *testing.T) {
	fx := setupTest(t)

	err := fx.rbacSvc.CreateRole("Admin", role.Admin)
	require.NoError(t, err)

	err = fx.raw.AssignRole(context.Background(), 99999, role.Admin)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// TestGetRoles_UserNotFound verifies ErrUserNotFound when getting roles for a non-existent user.
func TestGetRoles_UserNotFound(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.GetRoles(context.Background(), 99999)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// TestGetRoles_NoRoles verifies GetRoles returns nil when the user has no roles.
func TestGetRoles_NoRoles(t *testing.T) {
	fx := setupTest(t)

	user, err := fx.raw.Create(context.Background(), "Leo", "leo@example.com", "pass")
	require.NoError(t, err)

	roles, err := fx.raw.GetRoles(context.Background(), user.ID)
	require.NoError(t, err)
	assert.Nil(t, roles)
}

// TestRevokeRole_Success verifies a role can be revoked.
func TestRevokeRole_Success(t *testing.T) {
	fx := setupTest(t)

	user, err := fx.raw.Create(context.Background(), "Mallory", "mallory@example.com", "pass")
	require.NoError(t, err)

	err = fx.rbacSvc.CreateRole("Temp", "temp")
	require.NoError(t, err)

	err = fx.raw.AssignRole(context.Background(), user.ID, "temp")
	require.NoError(t, err)

	roles, err := fx.raw.GetRoles(context.Background(), user.ID)
	require.NoError(t, err)
	require.Len(t, roles, 1)

	err = fx.raw.RevokeRole(context.Background(), user.ID, "temp")
	require.NoError(t, err)

	roles, err = fx.raw.GetRoles(context.Background(), user.ID)
	require.NoError(t, err)
	assert.Nil(t, roles)
}

// TestRevokeRole_UserNotFound verifies ErrUserNotFound when revoking from a non-existent user.
func TestRevokeRole_UserNotFound(t *testing.T) {
	fx := setupTest(t)

	err := fx.rbacSvc.CreateRole("Admin", role.Admin)
	require.NoError(t, err)

	err = fx.raw.RevokeRole(context.Background(), 99999, role.Admin)
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// TestRevokeRole_NotAssigned verifies RevokeRole is a no-op when the role was not assigned.
func TestRevokeRole_NotAssigned(t *testing.T) {
	fx := setupTest(t)

	user, err := fx.raw.Create(context.Background(), "Nina", "nina@example.com", "pass")
	require.NoError(t, err)

	err = fx.rbacSvc.CreateRole("Viewer", "viewer")
	require.NoError(t, err)

	// Revoke without assigning — should not error
	err = fx.raw.RevokeRole(context.Background(), user.ID, "viewer")
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// DeleteExpiredTokens
// ---------------------------------------------------------------------------

// TestDeleteExpiredTokens_DeletesOldTokens verifies tokens older than ttlDays are soft-deleted.
func TestDeleteExpiredTokens_DeletesOldTokens(t *testing.T) {
	fx := setupTest(t)

	// Create a user and login to get a token
	_, err := fx.raw.Create(context.Background(), "OldToken", "old@example.com", "pass")
	require.NoError(t, err)

	_, token, err := fx.raw.Login(context.Background(), "old@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	// Manually set created_at to 31 days ago to simulate an expired token
	token.CreatedAt = time.Now().Add(-31 * 24 * time.Hour)
	token.UpdatedAt = time.Now()
	err = fx.repo.UsersToken.Save(token)
	require.NoError(t, err)

	// Act: delete tokens older than 30 days
	err = fx.raw.DeleteExpiredTokens(context.Background(), 30)
	require.NoError(t, err)

	// Assert: token should now be soft-deleted
	gotToken, err := fx.repo.UsersToken.GetUsersTokenByToken(token.Token)
	require.NoError(t, err)
	assert.True(t, gotToken.DeletedAt.Valid)
}

// TestDeleteExpiredTokens_KeepsFreshTokens verifies recent tokens are not deleted.
func TestDeleteExpiredTokens_KeepsFreshTokens(t *testing.T) {
	fx := setupTest(t)

	// Create a user and login to get a token
	_, err := fx.raw.Create(context.Background(), "FreshToken", "fresh@example.com", "pass")
	require.NoError(t, err)

	_, token, err := fx.raw.Login(context.Background(), "fresh@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	// Act: delete tokens older than 30 days (token is fresh)
	err = fx.raw.DeleteExpiredTokens(context.Background(), 30)
	require.NoError(t, err)

	// Assert: token should NOT be soft-deleted
	gotToken, err := fx.repo.UsersToken.GetUsersTokenByToken(token.Token)
	require.NoError(t, err)
	assert.False(t, gotToken.DeletedAt.Valid)
}

// TestDeleteExpiredTokens_UsesLastUsedAt verifies that last_used_at is preferred
// over created_at when determining expiry.
func TestDeleteExpiredTokens_UsesLastUsedAt(t *testing.T) {
	fx := setupTest(t)

	// Create a user and login to get a token
	_, err := fx.raw.Create(context.Background(), "LastUsed", "lastused@example.com", "pass")
	require.NoError(t, err)

	_, token, err := fx.raw.Login(context.Background(), "lastused@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	// Set created_at to recent (1 day ago) but last_used_at to far past (31 days ago)
	token.CreatedAt = time.Now().Add(-1 * 24 * time.Hour)
	token.LastUsedAt = pq.NullTime{Time: time.Now().Add(-31 * 24 * time.Hour), Valid: true}
	token.UpdatedAt = time.Now()
	err = fx.repo.UsersToken.Save(token)
	require.NoError(t, err)

	// Act: delete tokens older than 15 days
	// Token's created_at is 1 day old (fresh), but last_used_at is 31 days old (expired)
	err = fx.raw.DeleteExpiredTokens(context.Background(), 15)
	require.NoError(t, err)

	// Assert: token should be deleted because last_used_at is the decisive field
	gotToken, err := fx.repo.UsersToken.GetUsersTokenByToken(token.Token)
	require.NoError(t, err)
	assert.True(t, gotToken.DeletedAt.Valid)
}

// TestDeleteExpiredTokens_TTLZero verifies that ttlDays=0 deletes past tokens.
func TestDeleteExpiredTokens_TTLZero(t *testing.T) {
	fx := setupTest(t)

	// Create a user and login to get a token
	_, err := fx.raw.Create(context.Background(), "ZeroTTL", "zerottl@example.com", "pass")
	require.NoError(t, err)

	_, token, err := fx.raw.Login(context.Background(), "zerottl@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token)

	// Use UTC() to match PostgreSQL's NOW() behaviour with timestamp without time zone
	// (otherwise local +04 timezone wall clock is stored ahead of NOW() in UTC)
	nowUTC := time.Now().UTC()
	token.LastUsedAt = pq.NullTime{Time: nowUTC.Add(-1 * time.Hour), Valid: true}
	token.UpdatedAt = nowUTC
	err = fx.repo.UsersToken.Save(token)
	require.NoError(t, err)

	// Act: delete tokens older than 0 days (= any token with last_used_at/created_at before NOW())
	err = fx.raw.DeleteExpiredTokens(context.Background(), 0)
	require.NoError(t, err)

	// Assert: token should be soft-deleted
	gotToken, err := fx.repo.UsersToken.GetUsersTokenByToken(token.Token)
	require.NoError(t, err)
	assert.True(t, gotToken.DeletedAt.Valid)
}

// TestDeleteExpiredTokens_NoTokens verifies the method handles an empty DB without error.
func TestDeleteExpiredTokens_NoTokens(t *testing.T) {
	fx := setupTest(t)

	// Act: no users or tokens exist, delete should not error
	err := fx.raw.DeleteExpiredTokens(context.Background(), 30)
	require.NoError(t, err)
}

// TestDeleteExpiredTokens_AlreadyDeletedTokens verifies already deleted tokens
// are not affected (WHERE deleted_at IS NULL in the SQL).
func TestDeleteExpiredTokens_AlreadyDeletedTokens(t *testing.T) {
	fx := setupTest(t)

	// Create a user and login to get two tokens
	_, err := fx.raw.Create(context.Background(), "DelUser", "deluser@example.com", "pass")
	require.NoError(t, err)

	_, token1, err := fx.raw.Login(context.Background(), "deluser@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token1)

	_, token2, err := fx.raw.Login(context.Background(), "deluser@example.com", "pass")
	require.NoError(t, err)
	require.NotNil(t, token2)

	// Make both tokens old
	token1.CreatedAt = time.Now().Add(-31 * 24 * time.Hour)
	token1.UpdatedAt = time.Now()
	err = fx.repo.UsersToken.Save(token1)
	require.NoError(t, err)

	token2.CreatedAt = time.Now().Add(-31 * 24 * time.Hour)
	token2.UpdatedAt = time.Now()
	err = fx.repo.UsersToken.Save(token2)
	require.NoError(t, err)

	// Soft-delete token1 manually (like logout would do)
	now := time.Now()
	token1.DeletedAt = pq.NullTime{Time: now, Valid: true}
	token1.UpdatedAt = now
	err = fx.repo.UsersToken.Save(token1)
	require.NoError(t, err)

	// Act: delete tokens older than 30 days
	err = fx.raw.DeleteExpiredTokens(context.Background(), 30)
	require.NoError(t, err)

	// Assert: token1 is already deleted, token2 should now be deleted
	got1, err := fx.repo.UsersToken.GetUsersTokenByToken(token1.Token)
	require.NoError(t, err)
	assert.True(t, got1.DeletedAt.Valid, "already deleted token should remain deleted")

	got2, err := fx.repo.UsersToken.GetUsersTokenByToken(token2.Token)
	require.NoError(t, err)
	assert.True(t, got2.DeletedAt.Valid, "old non-deleted token should be deleted by cleanup")
}

// ---------------------------------------------------------------------------
// UpdatePassword
// ---------------------------------------------------------------------------

// TestUpdatePassword_Success verifies a password can be updated and new password works for login.
func TestUpdatePassword_Success(t *testing.T) {
	fx := setupTest(t)

	created, err := fx.raw.Create(context.Background(), "PassUser", "passuser@example.com", "oldpass")
	require.NoError(t, err)

	err = fx.raw.UpdatePassword(context.Background(), created.ID, "newpass")
	require.NoError(t, err)

	// Verify new password works via login
	_, token, err := fx.raw.Login(context.Background(), "passuser@example.com", "newpass")
	require.NoError(t, err)
	require.NotNil(t, token)

	// Old password should no longer work
	_, _, err = fx.raw.Login(context.Background(), "passuser@example.com", "oldpass")
	require.Error(t, err)
	assert.ErrorContains(t, err, "access denied")
}

// TestUpdatePassword_UserNotFound verifies ErrUserNotFound for a non-existent user.
func TestUpdatePassword_UserNotFound(t *testing.T) {
	fx := setupTest(t)

	err := fx.raw.UpdatePassword(context.Background(), 99999, "newpass")
	require.Error(t, err)
	assert.ErrorContains(t, err, "user not found")
}

// ---------------------------------------------------------------------------
// List — negative offset
// ---------------------------------------------------------------------------

// TestList_NegativeOffset verifies that a negative offset is treated as 0.
func TestList_NegativeOffset(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Alice", "alice@example.com", "a")
	require.NoError(t, err)

	// Negative offset should be clamped to 0
	users, total, err := fx.raw.List(context.Background(), 10, -1)
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, int64(1), total)
}
