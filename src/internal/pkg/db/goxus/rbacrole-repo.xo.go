// Code generated from rbacrole.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// IRbacRoleRepository defines the repository interface
type IRbacRoleRepository interface {
	Save(rr *RbacRole) error
	Delete(rr *RbacRole) error
	GetAll() ([]*RbacRole, error)
	GetLastID() (*RbacRole, error)
	GetRbacRoleByID(id int64) (*RbacRole, error)
	GetRbacRoleBySlug(slug string) (*RbacRole, error)
	GetRolesByUserID(userID int64) ([]*RbacRole, error)
}

// Save saves the RbacRole to the database.
func (repo *RbacRoleRepository) Save(rr *RbacRole) error {
	return rr.Save(repo.db)
}

// Delete deletes the RbacRole from the database.
func (repo *RbacRoleRepository) Delete(rr *RbacRole) error {
	return rr.Delete(repo.db)
}

// RbacRoleRepository реализует работу с таблицей 'rbac_roles'.
type RbacRoleRepository struct {
	db pgxdb.DBQuery
}

// NewRbacRoleRepository создает новый репозиторий.
func NewRbacRoleRepository(db pgxdb.DBQuery) *RbacRoleRepository {
	return &RbacRoleRepository{db: db}
}

// GetAll возвращает все записи
func (repo *RbacRoleRepository) GetAll() ([]*RbacRole, error) {
	return GetAllRbacRole(repo.db)
}

// GetLastID возвращает последний ID
func (repo *RbacRoleRepository) GetLastID() (*RbacRole, error) {
	return GetLastIDRbacRole(repo.db)
}

// GetRbacRoleByID возвращает одну запись по индексу 'rbac_roles_pk'.
func (repo *RbacRoleRepository) GetRbacRoleByID(id int64) (*RbacRole, error) {
	return GetRbacRoleByID(repo.db, id)
}

// GetRbacRoleBySlug возвращает одну запись по индексу 'rbac_roles_slug_uindex'.
func (repo *RbacRoleRepository) GetRbacRoleBySlug(slug string) (*RbacRole, error) {
	return GetRbacRoleBySlug(repo.db, slug)
}

// GetRolesByUserID runs a custom query, returning results as RbacRole.
func (repo *RbacRoleRepository) GetRolesByUserID(userID int64) ([]*RbacRole, error) {
	return GetRolesByUserID(repo.db, userID)
}
