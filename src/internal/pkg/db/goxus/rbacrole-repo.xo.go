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
	GetAllWithPagination(limit, offset int) ([]*RbacRole, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*RbacRole, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*RbacRole, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*RbacRole, error)
	GetRbacRoleByID(id int64) (*RbacRole, error)
	GetRbacRoleByIDCount(id int64) (int64, error)
	GetRbacRoleBySlug(slug string) (*RbacRole, error)
	GetRbacRoleBySlugCount(slug string) (int64, error)
	GetRolesByUserID(userID int64) ([]*RbacRole, error)
	GetRolesByUserIDCount(userID int64) (int64, error)
	GetRolesByUserIDWithPagination(userID int64, limit, offset int) ([]*RbacRole, error)
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

// GetAllWithPagination возвращает записи с пагинацией
func (repo *RbacRoleRepository) GetAllWithPagination(limit, offset int) ([]*RbacRole, error) {
	return GetAllRbacRoleWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *RbacRoleRepository) GetAllCount() (int64, error) {
	return GetAllRbacRoleCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *RbacRoleRepository) GetBySQL(sqlstr string, args ...any) ([]*RbacRole, error) {
	return GetRbacRolesBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *RbacRoleRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*RbacRole, error) {
	return GetRbacRolesBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *RbacRoleRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetRbacRolesBySQLCount(repo.db, sqlstr, args...)
}

// GetLastID возвращает последний ID
func (repo *RbacRoleRepository) GetLastID() (*RbacRole, error) {
	return GetLastIDRbacRole(repo.db)
}

// GetRbacRoleByID возвращает одну запись по индексу 'rbac_roles_pk'.
func (repo *RbacRoleRepository) GetRbacRoleByID(id int64) (*RbacRole, error) {
	return GetRbacRoleByID(repo.db, id)
}

// GetRbacRoleByIDCount возвращает количество записей по индексу 'rbac_roles_pk'.
func (repo *RbacRoleRepository) GetRbacRoleByIDCount(id int64) (int64, error) {
	return GetRbacRoleByIDCount(repo.db, id)
}

// GetRbacRoleBySlug возвращает одну запись по индексу 'rbac_roles_slug_uindex'.
func (repo *RbacRoleRepository) GetRbacRoleBySlug(slug string) (*RbacRole, error) {
	return GetRbacRoleBySlug(repo.db, slug)
}

// GetRbacRoleBySlugCount возвращает количество записей по индексу 'rbac_roles_slug_uindex'.
func (repo *RbacRoleRepository) GetRbacRoleBySlugCount(slug string) (int64, error) {
	return GetRbacRoleBySlugCount(repo.db, slug)
}

// GetRolesByUserID runs a custom query, returning results as RbacRole.
func (repo *RbacRoleRepository) GetRolesByUserID(userID int64) ([]*RbacRole, error) {
	return GetRolesByUserID(repo.db, userID)
}

// GetRolesByUserIDCount runs a custom count query from repository
func (repo *RbacRoleRepository) GetRolesByUserIDCount(userID int64) (int64, error) {
	return GetRolesByUserIDCount(repo.db, userID)
}

// GetRolesByUserIDWithPagination runs a custom query with pagination from repository
func (repo *RbacRoleRepository) GetRolesByUserIDWithPagination(userID int64, limit, offset int) ([]*RbacRole, error) {
	return GetRolesByUserIDWithPagination(repo.db, userID, limit, offset)
}
