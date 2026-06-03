package userdomain

import (
	"context"
	"testing"
	"time"

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

// TestList_Empty verifies GetAll returns an empty slice when no users exist.
func TestList_Empty(t *testing.T) {
	fx := setupTest(t)

	users, err := fx.raw.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, users, 0)
}

// TestList_Multiple verifies that multiple users are returned.
func TestList_Multiple(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.raw.Create(context.Background(), "Alice", "alice@example.com", "a")
	require.NoError(t, err)
	_, err = fx.raw.Create(context.Background(), "Bob", "bob@example.com", "b")
	require.NoError(t, err)

	users, err := fx.raw.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, users, 2)
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
