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
	GetAllWithPagination(limit, offset int) ([]*RbacPermission, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*RbacPermission, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*RbacPermission, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*RbacPermission, error)
	GetRbacPermissionByID(id int64) (*RbacPermission, error)
	GetRbacPermissionBySlug(slug string) (*RbacPermission, error)
	GetPermissionsByRoleSlug(roleSlug string) ([]*RbacPermission, error)
	GetPermissionsByRoleSlugCount(roleSlug string) (int64, error)
	GetPermissionsByRoleSlugWithPagination(roleSlug string, limit, offset int) ([]*RbacPermission, error)
	GetPermissionsByUserIDAndSlug(userID int64, permSlug string) ([]*RbacPermission, error)
	GetPermissionsByUserIDAndSlugCount(userID int64, permSlug string) (int64, error)
	GetPermissionsByUserIDAndSlugWithPagination(userID int64, permSlug string, limit, offset int) ([]*RbacPermission, error)
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

// GetAllWithPagination возвращает записи с пагинацией
func (repo *RbacPermissionRepository) GetAllWithPagination(limit, offset int) ([]*RbacPermission, error) {
	return GetAllRbacPermissionWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *RbacPermissionRepository) GetAllCount() (int64, error) {
	return GetAllRbacPermissionCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *RbacPermissionRepository) GetBySQL(sqlstr string, args ...any) ([]*RbacPermission, error) {
	return GetRbacPermissionsBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *RbacPermissionRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*RbacPermission, error) {
	return GetRbacPermissionsBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *RbacPermissionRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetRbacPermissionsBySQLCount(repo.db, sqlstr, args...)
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

// GetPermissionsByRoleSlug runs a custom query, returning results as RbacPermission.
func (repo *RbacPermissionRepository) GetPermissionsByRoleSlug(roleSlug string) ([]*RbacPermission, error) {
	return GetPermissionsByRoleSlug(repo.db, roleSlug)
}

// GetPermissionsByRoleSlugCount runs a custom count query from repository
func (repo *RbacPermissionRepository) GetPermissionsByRoleSlugCount(roleSlug string) (int64, error) {
	return GetPermissionsByRoleSlugCount(repo.db, roleSlug)
}

// GetPermissionsByRoleSlugWithPagination runs a custom query with pagination from repository
func (repo *RbacPermissionRepository) GetPermissionsByRoleSlugWithPagination(roleSlug string, limit, offset int) ([]*RbacPermission, error) {
	return GetPermissionsByRoleSlugWithPagination(repo.db, roleSlug, limit, offset)
}

// GetPermissionsByUserIDAndSlug runs a custom query, returning results as RbacPermission.
func (repo *RbacPermissionRepository) GetPermissionsByUserIDAndSlug(userID int64, permSlug string) ([]*RbacPermission, error) {
	return GetPermissionsByUserIDAndSlug(repo.db, userID, permSlug)
}

// GetPermissionsByUserIDAndSlugCount runs a custom count query from repository
func (repo *RbacPermissionRepository) GetPermissionsByUserIDAndSlugCount(userID int64, permSlug string) (int64, error) {
	return GetPermissionsByUserIDAndSlugCount(repo.db, userID, permSlug)
}

// GetPermissionsByUserIDAndSlugWithPagination runs a custom query with pagination from repository
func (repo *RbacPermissionRepository) GetPermissionsByUserIDAndSlugWithPagination(userID int64, permSlug string, limit, offset int) ([]*RbacPermission, error) {
	return GetPermissionsByUserIDAndSlugWithPagination(repo.db, userID, permSlug, limit, offset)
}
