package rbac

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateRole_Success verifies a role is created and returned via GetAllRoles.
func TestCreateRole_Success(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Admin", "admin")
	require.NoError(t, err)

	roles, err := fx.svc.GetAllRoles()
	require.NoError(t, err)
	require.Len(t, roles, 1)
	assert.Equal(t, "Admin", roles[0].Name)
	assert.Equal(t, "admin", roles[0].Slug)
}

// TestCreateRole_Duplicate verifies that creating a role with an existing slug returns ErrAlreadyExists.
func TestCreateRole_Duplicate(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Admin", "admin")
	require.NoError(t, err)

	err = fx.svc.CreateRole("Admin Duplicate", "admin")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAlreadyExists)
}

// TestCreatePermission_Success verifies a permission is created and returned via GetAllPermissions.
func TestCreatePermission_Success(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreatePermission("Read", "read")
	require.NoError(t, err)

	perms, err := fx.svc.GetAllPermissions()
	require.NoError(t, err)
	require.Len(t, perms, 1)
	assert.Equal(t, "Read", perms[0].Name)
	assert.Equal(t, "read", perms[0].Slug)
}

// TestCreatePermission_Duplicate verifies that creating a permission with an existing slug returns ErrAlreadyExists.
func TestCreatePermission_Duplicate(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreatePermission("Read", "read")
	require.NoError(t, err)

	err = fx.svc.CreatePermission("Read Duplicate", "read")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAlreadyExists)
}

// TestAssignPermissionsToRole verifies multiple permissions can be assigned to a role
// and checked individually via CheckRolePermission.
func TestAssignPermissionsToRole(t *testing.T) {
	fx := setupTest(t)

	// Create role and permissions
	err := fx.svc.CreateRole("Editor", "editor")
	require.NoError(t, err)
	err = fx.svc.CreatePermission("Write", "write")
	require.NoError(t, err)
	err = fx.svc.CreatePermission("Read", "read")
	require.NoError(t, err)

	// Assign permissions to role
	err = fx.svc.AssignPermissionsToRole("editor", []string{"write", "read"})
	require.NoError(t, err)

	// Check role permissions
	perms, err := fx.svc.GetRolePermissions("editor")
	require.NoError(t, err)
	require.Len(t, perms, 2)

	permSlugs := make(map[string]bool)
	for _, p := range perms {
		permSlugs[p.Slug] = true
	}
	assert.True(t, permSlugs["write"])
	assert.True(t, permSlugs["read"])

	// Check each permission individually
	hasWrite, err := fx.svc.CheckRolePermission("editor", "write")
	require.NoError(t, err)
	assert.True(t, hasWrite)

	hasRead, err := fx.svc.CheckRolePermission("editor", "read")
	require.NoError(t, err)
	assert.True(t, hasRead)
}

// TestAssignPermissionsToRole_InvalidRole verifies ErrRoleNotFound for a non-existent role.
func TestAssignPermissionsToRole_InvalidRole(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.AssignPermissionsToRole("nonexistent", []string{"write"})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

// TestAssignPermissionsToRole_InvalidPermission verifies ErrPermissionNotFound for a non-existent permission.
func TestAssignPermissionsToRole_InvalidPermission(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Editor", "editor")
	require.NoError(t, err)

	err = fx.svc.AssignPermissionsToRole("editor", []string{"nonexistent"})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPermissionNotFound)
}

// TestAssignRoleToUser verifies a role is assigned to a user and returned via GetUserRoles.
func TestAssignRoleToUser(t *testing.T) {
	fx := setupTest(t)

	// Create role
	err := fx.svc.CreateRole("Moderator", "moderator")
	require.NoError(t, err)

	// Assign to test user
	err = fx.svc.AssignRoleToUser(testUserID, "moderator")
	require.NoError(t, err)

	// Check user role
	hasRole, err := fx.svc.CheckUserRole(testUserID, "moderator")
	require.NoError(t, err)
	assert.True(t, hasRole)

	// Get user roles
	roles, err := fx.svc.GetUserRoles(testUserID)
	require.NoError(t, err)
	require.Len(t, roles, 1)
	assert.Equal(t, "Moderator", roles[0].Name)
}

// TestAssignRoleToUser_InvalidRole verifies ErrRoleNotFound for a non-existent role.
func TestAssignRoleToUser_InvalidRole(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.AssignRoleToUser(testUserID, "nonexistent")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

// TestAssignRoleToUser_Duplicate verifies ErrAlreadyExists when assigning the same role twice.
func TestAssignRoleToUser_Duplicate(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Moderator", "moderator")
	require.NoError(t, err)

	err = fx.svc.AssignRoleToUser(testUserID, "moderator")
	require.NoError(t, err)

	err = fx.svc.AssignRoleToUser(testUserID, "moderator")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAlreadyExists)
}

// TestCheckUserRole_NoRole verifies CheckUserRole returns false when the role exists but is not assigned to the user.
func TestCheckUserRole_NoRole(t *testing.T) {
	fx := setupTest(t)

	// Create a role but don't assign it to the user
	err := fx.svc.CreateRole("Admin", "admin")
	require.NoError(t, err)

	hasRole, err := fx.svc.CheckUserRole(testUserID, "admin")
	require.NoError(t, err)
	assert.False(t, hasRole)
}

// TestCheckUserPermission verifies CheckUserPermission returns true when the user has a role with that permission.
func TestCheckUserPermission(t *testing.T) {
	fx := setupTest(t)

	// Setup: role with permission assigned to user
	err := fx.svc.CreateRole("Viewer", "viewer")
	require.NoError(t, err)
	err = fx.svc.CreatePermission("View", "view")
	require.NoError(t, err)
	err = fx.svc.AssignPermissionsToRole("viewer", []string{"view"})
	require.NoError(t, err)
	err = fx.svc.AssignRoleToUser(testUserID, "viewer")
	require.NoError(t, err)

	// Check user permission
	hasPerm, err := fx.svc.CheckUserPermission(testUserID, "view")
	require.NoError(t, err)
	assert.True(t, hasPerm)
}

// TestCheckUserPermission_NoPermission verifies CheckUserPermission returns false
// when no role or user has the permission assigned.
func TestCheckUserPermission_NoPermission(t *testing.T) {
	fx := setupTest(t)

	// Create a permission but don't assign it to any role or user
	err := fx.svc.CreatePermission("View", "view")
	require.NoError(t, err)

	hasPerm, err := fx.svc.CheckUserPermission(testUserID, "view")
	require.NoError(t, err)
	assert.False(t, hasPerm)
}

// TestCheckUserPermission_UserHasNoRole verifies CheckUserPermission returns false
// when the user has no roles at all, even if the permission exists.
func TestCheckUserPermission_UserHasNoRole(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreatePermission("View", "view")
	require.NoError(t, err)

	hasPerm, err := fx.svc.CheckUserPermission(testUserID, "view")
	require.NoError(t, err)
	assert.False(t, hasPerm)
}

// TestRevokeUserRole verifies that revoking a user role removes the assignment.
func TestRevokeUserRole(t *testing.T) {
	fx := setupTest(t)

	// Setup
	err := fx.svc.CreateRole("Temp", "temp")
	require.NoError(t, err)
	err = fx.svc.AssignRoleToUser(testUserID, "temp")
	require.NoError(t, err)

	// Verify assigned
	hasRole, err := fx.svc.CheckUserRole(testUserID, "temp")
	require.NoError(t, err)
	assert.True(t, hasRole)

	// Revoke
	err = fx.svc.RevokeUserRole(testUserID, "temp")
	require.NoError(t, err)

	// Verify revoked
	hasRole, err = fx.svc.CheckUserRole(testUserID, "temp")
	require.NoError(t, err)
	assert.False(t, hasRole)
}

// TestRevokeUserRole_NotAssigned verifies RevokeUserRole is a no-op when the role was not assigned.
func TestRevokeUserRole_NotAssigned(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Temp", "temp")
	require.NoError(t, err)

	// Revoke without assigning — should not error
	err = fx.svc.RevokeUserRole(testUserID, "temp")
	require.NoError(t, err)
}

// TestRevokeRolePermission verifies that revoking a role permission removes the link.
func TestRevokeRolePermission(t *testing.T) {
	fx := setupTest(t)

	// Setup
	err := fx.svc.CreateRole("Editor", "editor")
	require.NoError(t, err)
	err = fx.svc.CreatePermission("Write", "write")
	require.NoError(t, err)
	err = fx.svc.AssignPermissionsToRole("editor", []string{"write"})
	require.NoError(t, err)

	// Verify assigned
	hasPerm, err := fx.svc.CheckRolePermission("editor", "write")
	require.NoError(t, err)
	assert.True(t, hasPerm)

	// Revoke
	err = fx.svc.RevokeRolePermission("editor", "write")
	require.NoError(t, err)

	// Verify revoked
	hasPerm, err = fx.svc.CheckRolePermission("editor", "write")
	require.NoError(t, err)
	assert.False(t, hasPerm)
}

// TestDeleteRole_Success verifies a bare role (no permissions, no users) can be deleted.
func TestDeleteRole_Success(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Temp", "temp")
	require.NoError(t, err)

	err = fx.svc.DeleteRole("temp")
	require.NoError(t, err)

	roles, err := fx.svc.GetAllRoles()
	require.NoError(t, err)
	assert.Len(t, roles, 0)
}

// TestDeleteRole_NotFound verifies ErrRoleNotFound when deleting a non-existent role.
func TestDeleteRole_NotFound(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.DeleteRole("nonexistent")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

// TestDeleteRole_InUse verifies ErrRoleInUse when the role is assigned to a user.
func TestDeleteRole_InUse(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Admin", "admin")
	require.NoError(t, err)
	err = fx.svc.AssignRoleToUser(testUserID, "admin")
	require.NoError(t, err)

	err = fx.svc.DeleteRole("admin")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRoleInUse)
}

// TestDeletePermission_Success verifies a bare permission (no role links) can be deleted.
func TestDeletePermission_Success(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreatePermission("Temp", "temp")
	require.NoError(t, err)

	err = fx.svc.DeletePermission("temp")
	require.NoError(t, err)

	perms, err := fx.svc.GetAllPermissions()
	require.NoError(t, err)
	assert.Len(t, perms, 0)
}

// TestDeletePermission_NotFound verifies ErrPermissionNotFound when deleting a non-existent permission.
func TestDeletePermission_NotFound(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.DeletePermission("nonexistent")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPermissionNotFound)
}

// TestDeletePermission_InUse verifies ErrPermissionInUse when the permission is assigned to a role.
func TestDeletePermission_InUse(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Admin", "admin")
	require.NoError(t, err)
	err = fx.svc.CreatePermission("Write", "write")
	require.NoError(t, err)
	err = fx.svc.AssignPermissionsToRole("admin", []string{"write"})
	require.NoError(t, err)

	err = fx.svc.DeletePermission("write")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPermissionInUse)
}

// TestGetAllRoles_Empty verifies GetAllRoles returns an empty slice when no roles exist.
func TestGetAllRoles_Empty(t *testing.T) {
	fx := setupTest(t)

	roles, err := fx.svc.GetAllRoles()
	require.NoError(t, err)
	assert.Len(t, roles, 0)
}

// TestGetAllPermissions_Empty verifies GetAllPermissions returns an empty slice when no permissions exist.
func TestGetAllPermissions_Empty(t *testing.T) {
	fx := setupTest(t)

	perms, err := fx.svc.GetAllPermissions()
	require.NoError(t, err)
	assert.Len(t, perms, 0)
}

// TestGetUserRoles_NoRoles verifies GetUserRoles returns nil when the user has no roles.
func TestGetUserRoles_NoRoles(t *testing.T) {
	fx := setupTest(t)

	roles, err := fx.svc.GetUserRoles(testUserID)
	require.NoError(t, err)
	assert.Nil(t, roles)
}

// TestGetRolePermissions_NoPermissions verifies GetRolePermissions returns nil when the role has no permissions.
func TestGetRolePermissions_NoPermissions(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Empty", "empty")
	require.NoError(t, err)

	perms, err := fx.svc.GetRolePermissions("empty")
	require.NoError(t, err)
	assert.Nil(t, perms)
}

// TestGetRolePermissions_NotFound verifies ErrRoleNotFound for a non-existent role.
func TestGetRolePermissions_NotFound(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.svc.GetRolePermissions("nonexistent")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

// TestCheckRolePermission_NotFound verifies ErrRoleNotFound for a non-existent role.
func TestCheckRolePermission_NotFound(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.svc.CheckRolePermission("nonexistent", "read")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

// TestCheckUserRole_NotFound verifies ErrRoleNotFound for a non-existent role slug.
func TestCheckUserRole_NotFound(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.svc.CheckUserRole(testUserID, "nonexistent")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

// TestCheckUserPermission_NotFound verifies ErrPermissionNotFound for a non-existent permission slug.
func TestCheckUserPermission_NotFound(t *testing.T) {
	fx := setupTest(t)

	_, err := fx.svc.CheckUserPermission(testUserID, "nonexistent")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPermissionNotFound)
}

// TestAssignPermissionsToRole_Idempotent verifies that assigning the same permission twice
// is idempotent — the second call hits the continue branch and does not duplicate the link.
func TestAssignPermissionsToRole_Idempotent(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Editor", "editor")
	require.NoError(t, err)
	err = fx.svc.CreatePermission("Write", "write")
	require.NoError(t, err)

	// First assignment
	err = fx.svc.AssignPermissionsToRole("editor", []string{"write"})
	require.NoError(t, err)

	perms, err := fx.svc.GetRolePermissions("editor")
	require.NoError(t, err)
	require.Len(t, perms, 1)

	// Second assignment — hits the continue branch, must not duplicate
	err = fx.svc.AssignPermissionsToRole("editor", []string{"write"})
	require.NoError(t, err)

	perms, err = fx.svc.GetRolePermissions("editor")
	require.NoError(t, err)
	require.Len(t, perms, 1)
}

// TestDeleteRole_WithPermissions verifies that deleting a role with assigned permissions
// cleans up role-permission links, removes the role, but leaves permissions intact.
func TestDeleteRole_WithPermissions(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Editor", "editor")
	require.NoError(t, err)
	err = fx.svc.CreatePermission("Write", "write")
	require.NoError(t, err)
	err = fx.svc.CreatePermission("Read", "read")
	require.NoError(t, err)

	// Assign permissions
	err = fx.svc.AssignPermissionsToRole("editor", []string{"write", "read"})
	require.NoError(t, err)

	// Delete the role
	err = fx.svc.DeleteRole("editor")
	require.NoError(t, err)

	// Role is gone
	roles, err := fx.svc.GetAllRoles()
	require.NoError(t, err)
	assert.Len(t, roles, 0)

	// Permissions survive independently
	perms, err := fx.svc.GetAllPermissions()
	require.NoError(t, err)
	assert.Len(t, perms, 2)

	// Role-permission links are cleaned — querying via deleted role returns ErrRoleNotFound
	_, err = fx.svc.GetRolePermissions("editor")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRoleNotFound)
}

// TestRevokeRolePermission_NotAssigned verifies that revoking a permission that was never
// assigned to the role is a no-op and does not return an error.
func TestRevokeRolePermission_NotAssigned(t *testing.T) {
	fx := setupTest(t)

	err := fx.svc.CreateRole("Guest", "guest")
	require.NoError(t, err)
	err = fx.svc.CreatePermission("Write", "write")
	require.NoError(t, err)

	// Revoke a permission that was never assigned — no error
	err = fx.svc.RevokeRolePermission("guest", "write")
	require.NoError(t, err)
}
