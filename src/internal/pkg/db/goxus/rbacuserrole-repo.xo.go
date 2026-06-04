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
	GetAllWithPagination(limit, offset int) ([]*RbacUserRole, error)
	GetBySQL(sqlstr string, args ...any) ([]*RbacUserRole, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*RbacUserRole, error)
	GetLastID() (*RbacUserRole, error)
	GetRbacUserRoleByID(id int64) (*RbacUserRole, error)
	GetRbacUserRoleByUserIDRoleID(userID int64, roleID int64) (*RbacUserRole, error)
	GetRbacUserRoleByRoleID(roleID int64) ([]*RbacUserRole, error)
	GetRbacUserRoleByRoleIDWithPagination(roleID int64, limit, offset int) ([]*RbacUserRole, error)
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

// GetAllWithPagination возвращает записи с пагинацией
func (repo *RbacUserRoleRepository) GetAllWithPagination(limit, offset int) ([]*RbacUserRole, error) {
	return GetAllRbacUserRoleWithPagination(repo.db, limit, offset)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *RbacUserRoleRepository) GetBySQL(sqlstr string, args ...any) ([]*RbacUserRole, error) {
	return GetRbacUserRolesBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *RbacUserRoleRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*RbacUserRole, error) {
	return GetRbacUserRolesBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
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

// GetRbacUserRoleByRoleID runs a custom query, returning results as RbacUserRole.
func (repo *RbacUserRoleRepository) GetRbacUserRoleByRoleID(roleID int64) ([]*RbacUserRole, error) {
	return GetRbacUserRoleByRoleID(repo.db, roleID)
}

// GetRbacUserRoleByRoleIDWithPagination runs a custom query with pagination from repository
func (repo *RbacUserRoleRepository) GetRbacUserRoleByRoleIDWithPagination(roleID int64, limit, offset int) ([]*RbacUserRole, error) {
	return GetRbacUserRoleByRoleIDWithPagination(repo.db, roleID, limit, offset)
}
