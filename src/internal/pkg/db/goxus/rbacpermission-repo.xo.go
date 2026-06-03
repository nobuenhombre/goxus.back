// Code generated from rbacpermission.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// IRbacPermissionRepository defines the repository interface
type IRbacPermissionRepository interface {
	Save(rp *RbacPermission) error
	Delete(rp *RbacPermission) error
	GetAll() ([]*RbacPermission, error)
	GetLastID() (*RbacPermission, error)
	GetRbacPermissionByID(id int64) (*RbacPermission, error)
	GetRbacPermissionBySlug(slug string) (*RbacPermission, error)
}

// Save saves the RbacPermission to the database.
func (repo *RbacPermissionRepository) Save(rp *RbacPermission) error {
	return rp.Save(repo.db)
}

// Delete deletes the RbacPermission from the database.
func (repo *RbacPermissionRepository) Delete(rp *RbacPermission) error {
	return rp.Delete(repo.db)
}

// RbacPermissionRepository реализует работу с таблицей 'rbac_permissions'.
type RbacPermissionRepository struct {
	db pgxdb.DBQuery
}

// NewRbacPermissionRepository создает новый репозиторий.
func NewRbacPermissionRepository(db pgxdb.DBQuery) *RbacPermissionRepository {
	return &RbacPermissionRepository{db: db}
}

// GetAll возвращает все записи
func (repo *RbacPermissionRepository) GetAll() ([]*RbacPermission, error) {
	return GetAllRbacPermission(repo.db)
}

// GetLastID возвращает последний ID
func (repo *RbacPermissionRepository) GetLastID() (*RbacPermission, error) {
	return GetLastIDRbacPermission(repo.db)
}

// GetRbacPermissionByID возвращает одну запись по индексу 'rbac_permissions_pk'.
func (repo *RbacPermissionRepository) GetRbacPermissionByID(id int64) (*RbacPermission, error) {
	return GetRbacPermissionByID(repo.db, id)
}

// GetRbacPermissionBySlug возвращает одну запись по индексу 'rbac_permissions_slug_uindex'.
func (repo *RbacPermissionRepository) GetRbacPermissionBySlug(slug string) (*RbacPermission, error) {
	return GetRbacPermissionBySlug(repo.db, slug)
}
