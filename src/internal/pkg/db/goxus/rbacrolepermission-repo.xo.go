// Code generated from rbacrolepermission.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// IRbacRolePermissionRepository defines the repository interface
type IRbacRolePermissionRepository interface {
	Save(rrp *RbacRolePermission) error
	Delete(rrp *RbacRolePermission) error
	GetAll() ([]*RbacRolePermission, error)
	GetLastID() (*RbacRolePermission, error)
	GetRbacRolePermissionByID(id int64) (*RbacRolePermission, error)
	GetRbacRolePermissionByRoleIDPermissionID(roleID int64, permissionID int64) (*RbacRolePermission, error)
	GetRbacRolePermissionByPermissionID(permissionID int64) ([]*RbacRolePermission, error)
	GetRbacRolePermissionByRoleID(roleID int64) ([]*RbacRolePermission, error)
}

// Save saves the RbacRolePermission to the database.
func (repo *RbacRolePermissionRepository) Save(rrp *RbacRolePermission) error {
	return rrp.Save(repo.db)
}

// Delete deletes the RbacRolePermission from the database.
func (repo *RbacRolePermissionRepository) Delete(rrp *RbacRolePermission) error {
	return rrp.Delete(repo.db)
}

// RbacRolePermissionRepository реализует работу с таблицей 'rbac_role_permissions'.
type RbacRolePermissionRepository struct {
	db pgxdb.DBQuery
}

// NewRbacRolePermissionRepository создает новый репозиторий.
func NewRbacRolePermissionRepository(db pgxdb.DBQuery) *RbacRolePermissionRepository {
	return &RbacRolePermissionRepository{db: db}
}

// GetAll возвращает все записи
func (repo *RbacRolePermissionRepository) GetAll() ([]*RbacRolePermission, error) {
	return GetAllRbacRolePermission(repo.db)
}

// GetLastID возвращает последний ID
func (repo *RbacRolePermissionRepository) GetLastID() (*RbacRolePermission, error) {
	return GetLastIDRbacRolePermission(repo.db)
}

// GetRbacRolePermissionByID возвращает одну запись по индексу 'rbac_role_permissions_pk'.
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByID(id int64) (*RbacRolePermission, error) {
	return GetRbacRolePermissionByID(repo.db, id)
}

// GetRbacRolePermissionByRoleIDPermissionID возвращает одну запись по индексу 'rbac_role_permissions_unique'.
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByRoleIDPermissionID(roleID int64, permissionID int64) (*RbacRolePermission, error) {
	return GetRbacRolePermissionByRoleIDPermissionID(repo.db, roleID, permissionID)
}

// GetRbacRolePermissionByPermissionID runs a custom query, returning results as RbacRolePermission.
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByPermissionID(permissionID int64) ([]*RbacRolePermission, error) {
	return GetRbacRolePermissionByPermissionID(repo.db, permissionID)
}

// GetRbacRolePermissionByRoleID runs a custom query, returning results as RbacRolePermission.
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByRoleID(roleID int64) ([]*RbacRolePermission, error) {
	return GetRbacRolePermissionByRoleID(repo.db, roleID)
}
