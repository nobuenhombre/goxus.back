// Package rbac provides Role-Based Access Control (RBAC) service.
package rbac

import (
	"errors"
	"fmt"

	"github.com/nobuenhombre/suikat/pkg/ge"

	"goxus/src/internal/pkg/db/goxus"
)

var (
	ErrPermissionInUse    = errors.New("cannot delete assigned permission")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrRoleInUse          = errors.New("cannot delete assigned role")
	ErrRoleNotFound       = errors.New("role not found")
	ErrAlreadyExists      = errors.New("already exists")
)

// Service defines the RBAC operations.
type Service interface {
	// CreateRole creates a new role.
	CreateRole(name, slug string) error
	// CreatePermission creates a new permission.
	CreatePermission(name, slug string) error

	// AssignPermissionsToRole assigns a list of permissions to a role.
	AssignPermissionsToRole(roleSlug string, permSlugs []string) error
	// AssignRoleToUser assigns a role to a user.
	AssignRoleToUser(userID int64, roleSlug string) error

	// CheckUserRole checks if a user has a given role.
	CheckUserRole(userID int64, roleSlug string) (bool, error)
	// CheckUserPermission checks if a user has a given permission (via roles).
	CheckUserPermission(userID int64, permSlug string) (bool, error)
	// CheckRolePermission checks if a role has a given permission.
	CheckRolePermission(roleSlug string, permSlug string) (bool, error)

	// RevokeUserRole removes a role from a user.
	RevokeUserRole(userID int64, roleSlug string) error
	// RevokeRolePermission removes a permission from a role.
	RevokeRolePermission(roleSlug string, permSlug string) error

	// GetAllRoles returns all roles.
	GetAllRoles() ([]*goxus.RbacRole, error)
	// GetUserRoles returns all roles assigned to a user.
	GetUserRoles(userID int64) ([]*goxus.RbacRole, error)
	// GetRolePermissions returns all permissions assigned to a role.
	GetRolePermissions(roleSlug string) ([]*goxus.RbacPermission, error)
	// GetAllPermissions returns all permissions.
	GetAllPermissions() ([]*goxus.RbacPermission, error)

	// DeleteRole deletes a role (fails if assigned to any user).
	DeleteRole(roleSlug string) error
	// DeletePermission deletes a permission (fails if assigned to any role).
	DeletePermission(permSlug string) error
}

// impl is the concrete implementation of Service.
type impl struct {
	repo *goxus.DbGoxusRepo
}

// New creates a new RBAC service.
func New(dbRepo *goxus.DbGoxusRepo) Service {
	return &impl{repo: dbRepo}
}

// CreateRole creates a new role.
func (s *impl) CreateRole(name, slug string) error {
	// check if role already exists
	existing, err := s.repo.RbacRole.GetRbacRoleBySlug(slug)
	if err == nil && existing != nil {
		return ge.Pin(fmt.Errorf("role '%s': %w", slug, ErrAlreadyExists))
	}

	role := &goxus.RbacRole{
		Name: name,
		Slug: slug,
	}

	err = s.repo.RbacRole.Save(role)
	if err != nil {
		return ge.Pin(err)
	}
	return nil
}

// CreatePermission creates a new permission.
func (s *impl) CreatePermission(name, slug string) error {
	existing, err := s.repo.RbacPermission.GetRbacPermissionBySlug(slug)
	if err == nil && existing != nil {
		return ge.Pin(fmt.Errorf("permission '%s': %w", slug, ErrAlreadyExists))
	}

	perm := &goxus.RbacPermission{
		Name: name,
		Slug: slug,
	}

	err = s.repo.RbacPermission.Save(perm)
	if err != nil {
		return ge.Pin(err)
	}
	return nil
}

// AssignPermissionsToRole assigns a list of permissions to a role.
func (s *impl) AssignPermissionsToRole(roleSlug string, permSlugs []string) error {
	// find role
	role, err := s.repo.RbacRole.GetRbacRoleBySlug(roleSlug)
	if err != nil {
		return ge.Pin(fmt.Errorf("role '%s': %w", roleSlug, ErrRoleNotFound))
	}

	// find all permissions
	for _, permSlug := range permSlugs {
		perm, err := s.repo.RbacPermission.GetRbacPermissionBySlug(permSlug)
		if err != nil {
			return ge.Pin(fmt.Errorf("permission '%s': %w", permSlug, ErrPermissionNotFound))
		}

		// check if already assigned
		existing, err := s.repo.RbacRolePermission.
			GetRbacRolePermissionByRoleIDPermissionID(role.ID, perm.ID)
		if err == nil && existing != nil {
			continue // already assigned, skip
		}

		rp := &goxus.RbacRolePermission{
			RoleID:       role.ID,
			PermissionID: perm.ID,
		}

		err = s.repo.RbacRolePermission.Save(rp)
		if err != nil {
			return ge.Pin(err)
		}
	}

	return nil
}

// AssignRoleToUser assigns a role to a user.
func (s *impl) AssignRoleToUser(userID int64, roleSlug string) error {
	role, err := s.repo.RbacRole.GetRbacRoleBySlug(roleSlug)
	if err != nil {
		return ge.Pin(fmt.Errorf("role '%s': %w", roleSlug, ErrRoleNotFound))
	}

	// check if already assigned
	existing, err := s.repo.RbacUserRole.GetRbacUserRoleByUserIDRoleID(userID, role.ID)
	if err == nil && existing != nil {
		return ge.Pin(fmt.Errorf("role '%s' already assigned to user '%d': %w",
			roleSlug, userID, ErrAlreadyExists))
	}

	ur := &goxus.RbacUserRole{
		UserID: userID,
		RoleID: role.ID,
	}

	err = s.repo.RbacUserRole.Save(ur)
	if err != nil {
		return ge.Pin(err)
	}
	return nil
}

// CheckUserRole checks if a user has a given role.
func (s *impl) CheckUserRole(userID int64, roleSlug string) (bool, error) {
	role, err := s.repo.RbacRole.GetRbacRoleBySlug(roleSlug)
	if err != nil {
		return false, ge.Pin(fmt.Errorf("role '%s': %w", roleSlug, ErrRoleNotFound))
	}

	ur, err := s.repo.RbacUserRole.GetRbacUserRoleByUserIDRoleID(userID, role.ID)
	if err != nil || ur == nil {
		return false, nil
	}

	return true, nil
}

// CheckUserPermission checks if a user has a given permission (via roles).
func (s *impl) CheckUserPermission(userID int64, permSlug string) (bool, error) {
	// validate permission exists
	_, err := s.repo.RbacPermission.GetRbacPermissionBySlug(permSlug)
	if err != nil {
		return false, ge.Pin(fmt.Errorf("permission '%s': %w", permSlug, ErrPermissionNotFound))
	}

	perms, err := s.repo.RbacPermission.GetPermissionsByUserIDAndSlug(userID, permSlug)
	if err != nil {
		return false, ge.Pin(err)
	}

	return len(perms) > 0, nil
}

// CheckRolePermission checks if a role has a given permission.
func (s *impl) CheckRolePermission(roleSlug string, permSlug string) (bool, error) {
	role, err := s.repo.RbacRole.GetRbacRoleBySlug(roleSlug)
	if err != nil {
		return false, ge.Pin(fmt.Errorf("role '%s': %w", roleSlug, ErrRoleNotFound))
	}

	perm, err := s.repo.RbacPermission.GetRbacPermissionBySlug(permSlug)
	if err != nil {
		return false, ge.Pin(fmt.Errorf("permission '%s': %w", permSlug, ErrPermissionNotFound))
	}

	rp, err := s.repo.RbacRolePermission.
		GetRbacRolePermissionByRoleIDPermissionID(role.ID, perm.ID)
	if err != nil || rp == nil {
		return false, nil
	}

	return true, nil
}

// RevokeUserRole removes a role from a user.
func (s *impl) RevokeUserRole(userID int64, roleSlug string) error {
	role, err := s.repo.RbacRole.GetRbacRoleBySlug(roleSlug)
	if err != nil {
		return ge.Pin(fmt.Errorf("role '%s': %w", roleSlug, ErrRoleNotFound))
	}

	ur, err := s.repo.RbacUserRole.GetRbacUserRoleByUserIDRoleID(userID, role.ID)
	if err != nil || ur == nil {
		return nil // not assigned, nothing to revoke
	}

	err = s.repo.RbacUserRole.Delete(ur)
	if err != nil {
		return ge.Pin(err)
	}
	return nil
}

// RevokeRolePermission removes a permission from a role.
func (s *impl) RevokeRolePermission(roleSlug string, permSlug string) error {
	role, err := s.repo.RbacRole.GetRbacRoleBySlug(roleSlug)
	if err != nil {
		return ge.Pin(fmt.Errorf("role '%s': %w", roleSlug, ErrRoleNotFound))
	}

	perm, err := s.repo.RbacPermission.GetRbacPermissionBySlug(permSlug)
	if err != nil {
		return ge.Pin(fmt.Errorf("permission '%s': %w", permSlug, ErrPermissionNotFound))
	}

	rp, err := s.repo.RbacRolePermission.
		GetRbacRolePermissionByRoleIDPermissionID(role.ID, perm.ID)
	if err != nil || rp == nil {
		return nil // not assigned, nothing to revoke
	}

	err = s.repo.RbacRolePermission.Delete(rp)
	if err != nil {
		return ge.Pin(err)
	}
	return nil
}

// GetAllRoles returns all roles.
func (s *impl) GetAllRoles() ([]*goxus.RbacRole, error) {
	res, err := s.repo.RbacRole.GetAll()
	if err != nil {
		return nil, ge.Pin(err)
	}
	return res, nil
}

// GetUserRoles returns all roles assigned to a user.
func (s *impl) GetUserRoles(userID int64) ([]*goxus.RbacRole, error) {
	res, err := s.repo.RbacRole.GetRolesByUserID(userID)
	if err != nil {
		return nil, ge.Pin(err)
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

// GetRolePermissions returns all permissions assigned to a role.
func (s *impl) GetRolePermissions(roleSlug string) ([]*goxus.RbacPermission, error) {
	// validate role exists
	_, err := s.repo.RbacRole.GetRbacRoleBySlug(roleSlug)
	if err != nil {
		return nil, ge.Pin(fmt.Errorf("role '%s': %w", roleSlug, ErrRoleNotFound))
	}

	res, err := s.repo.RbacPermission.GetPermissionsByRoleSlug(roleSlug)
	if err != nil {
		return nil, ge.Pin(err)
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

// GetAllPermissions returns all permissions.
func (s *impl) GetAllPermissions() ([]*goxus.RbacPermission, error) {
	res, err := s.repo.RbacPermission.GetAll()
	if err != nil {
		return nil, ge.Pin(err)
	}
	return res, nil
}

// DeleteRole deletes a role (fails if assigned to any user).
func (s *impl) DeleteRole(roleSlug string) error {
	role, err := s.repo.RbacRole.GetRbacRoleBySlug(roleSlug)
	if err != nil {
		return ge.Pin(fmt.Errorf("role '%s': %w", roleSlug, ErrRoleNotFound))
	}

	// check if any user has this role
	userRoles, err := s.repo.RbacUserRole.GetRbacUserRoleByRoleID(role.ID)
	if err != nil {
		return ge.Pin(err)
	}

	if len(userRoles) > 0 {
		return ge.Pin(fmt.Errorf("role '%s': %w", roleSlug, ErrRoleInUse))
	}

	// delete all role-permission links first
	rolePerms, err := s.repo.RbacRolePermission.GetRbacRolePermissionByRoleID(role.ID)
	if err != nil {
		return ge.Pin(err)
	}

	for _, rp := range rolePerms {
		err = s.repo.RbacRolePermission.Delete(rp)
		if err != nil {
			return ge.Pin(err)
		}
	}

	err = s.repo.RbacRole.Delete(role)
	if err != nil {
		return ge.Pin(err)
	}
	return nil
}

// DeletePermission deletes a permission (fails if assigned to any role).
func (s *impl) DeletePermission(permSlug string) error {
	perm, err := s.repo.RbacPermission.GetRbacPermissionBySlug(permSlug)
	if err != nil {
		return ge.Pin(fmt.Errorf("permission '%s': %w", permSlug, ErrPermissionNotFound))
	}

	// check if any role uses this permission
	rolePerms, err := s.repo.RbacRolePermission.GetRbacRolePermissionByPermissionID(perm.ID)
	if err != nil {
		return ge.Pin(err)
	}

	if len(rolePerms) > 0 {
		return ge.Pin(fmt.Errorf("permission '%s': %w", permSlug, ErrPermissionInUse))
	}

	err = s.repo.RbacPermission.Delete(perm)
	if err != nil {
		return ge.Pin(err)
	}
	return nil
}
