// Code generated from rbacuserrole.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// IRbacUserRoleRepository defines the repository interface
type IRbacUserRoleRepository interface {
	Save(rur *RbacUserRole) error
	Delete(rur *RbacUserRole) error
	GetAll() ([]*RbacUserRole, error)
	GetLastID() (*RbacUserRole, error)
	GetRbacUserRoleByID(id int64) (*RbacUserRole, error)
	GetRbacUserRoleByUserIDRoleID(userID int64, roleID int64) (*RbacUserRole, error)
}

// Save saves the RbacUserRole to the database.
func (repo *RbacUserRoleRepository) Save(rur *RbacUserRole) error {
	return rur.Save(repo.db)
}

// Delete deletes the RbacUserRole from the database.
func (repo *RbacUserRoleRepository) Delete(rur *RbacUserRole) error {
	return rur.Delete(repo.db)
}

// RbacUserRoleRepository реализует работу с таблицей 'rbac_user_roles'.
type RbacUserRoleRepository struct {
	db pgxdb.DBQuery
}

// NewRbacUserRoleRepository создает новый репозиторий.
func NewRbacUserRoleRepository(db pgxdb.DBQuery) *RbacUserRoleRepository {
	return &RbacUserRoleRepository{db: db}
}

// GetAll возвращает все записи
func (repo *RbacUserRoleRepository) GetAll() ([]*RbacUserRole, error) {
	return GetAllRbacUserRole(repo.db)
}

// GetLastID возвращает последний ID
func (repo *RbacUserRoleRepository) GetLastID() (*RbacUserRole, error) {
	return GetLastIDRbacUserRole(repo.db)
}

// GetRbacUserRoleByID возвращает одну запись по индексу 'rbac_user_roles_pk'.
func (repo *RbacUserRoleRepository) GetRbacUserRoleByID(id int64) (*RbacUserRole, error) {
	return GetRbacUserRoleByID(repo.db, id)
}

// GetRbacUserRoleByUserIDRoleID возвращает одну запись по индексу 'rbac_user_roles_unique'.
func (repo *RbacUserRoleRepository) GetRbacUserRoleByUserIDRoleID(userID int64, roleID int64) (*RbacUserRole, error) {
	return GetRbacUserRoleByUserIDRoleID(repo.db, userID, roleID)
}
