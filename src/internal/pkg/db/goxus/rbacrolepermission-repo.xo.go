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
	GetAllWithPagination(limit, offset int) ([]*RbacRolePermission, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*RbacRolePermission, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*RbacRolePermission, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*RbacRolePermission, error)
	GetRbacRolePermissionByID(id int64) (*RbacRolePermission, error)
	GetRbacRolePermissionByIDCount(id int64) (int64, error)
	GetRbacRolePermissionByRoleIDPermissionID(roleID int64, permissionID int64) (*RbacRolePermission, error)
	GetRbacRolePermissionByRoleIDPermissionIDCount(roleID int64, permissionID int64) (int64, error)
	GetRbacRolePermissionByPermissionID(permissionID int64) ([]*RbacRolePermission, error)
	GetRbacRolePermissionByPermissionIDCount(permissionID int64) (int64, error)
	GetRbacRolePermissionByPermissionIDWithPagination(permissionID int64, limit, offset int) ([]*RbacRolePermission, error)
	GetRbacRolePermissionByRoleID(roleID int64) ([]*RbacRolePermission, error)
	GetRbacRolePermissionByRoleIDCount(roleID int64) (int64, error)
	GetRbacRolePermissionByRoleIDWithPagination(roleID int64, limit, offset int) ([]*RbacRolePermission, error)
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

// GetAllWithPagination возвращает записи с пагинацией
func (repo *RbacRolePermissionRepository) GetAllWithPagination(limit, offset int) ([]*RbacRolePermission, error) {
	return GetAllRbacRolePermissionWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *RbacRolePermissionRepository) GetAllCount() (int64, error) {
	return GetAllRbacRolePermissionCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *RbacRolePermissionRepository) GetBySQL(sqlstr string, args ...any) ([]*RbacRolePermission, error) {
	return GetRbacRolePermissionsBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *RbacRolePermissionRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*RbacRolePermission, error) {
	return GetRbacRolePermissionsBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *RbacRolePermissionRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetRbacRolePermissionsBySQLCount(repo.db, sqlstr, args...)
}

// GetLastID возвращает последний ID
func (repo *RbacRolePermissionRepository) GetLastID() (*RbacRolePermission, error) {
	return GetLastIDRbacRolePermission(repo.db)
}

// GetRbacRolePermissionByID возвращает одну запись по индексу 'rbac_role_permissions_pk'.
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByID(id int64) (*RbacRolePermission, error) {
	return GetRbacRolePermissionByID(repo.db, id)
}

// GetRbacRolePermissionByIDCount возвращает количество записей по индексу 'rbac_role_permissions_pk'.
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByIDCount(id int64) (int64, error) {
	return GetRbacRolePermissionByIDCount(repo.db, id)
}

// GetRbacRolePermissionByRoleIDPermissionID возвращает одну запись по индексу 'rbac_role_permissions_unique'.
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByRoleIDPermissionID(roleID int64, permissionID int64) (*RbacRolePermission, error) {
	return GetRbacRolePermissionByRoleIDPermissionID(repo.db, roleID, permissionID)
}

// GetRbacRolePermissionByRoleIDPermissionIDCount возвращает количество записей по индексу 'rbac_role_permissions_unique'.
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByRoleIDPermissionIDCount(roleID int64, permissionID int64) (int64, error) {
	return GetRbacRolePermissionByRoleIDPermissionIDCount(repo.db, roleID, permissionID)
}

// GetRbacRolePermissionByPermissionID runs a custom query, returning results as RbacRolePermission.
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByPermissionID(permissionID int64) ([]*RbacRolePermission, error) {
	return GetRbacRolePermissionByPermissionID(repo.db, permissionID)
}

// GetRbacRolePermissionByPermissionIDCount runs a custom count query from repository
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByPermissionIDCount(permissionID int64) (int64, error) {
	return GetRbacRolePermissionByPermissionIDCount(repo.db, permissionID)
}

// GetRbacRolePermissionByPermissionIDWithPagination runs a custom query with pagination from repository
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByPermissionIDWithPagination(permissionID int64, limit, offset int) ([]*RbacRolePermission, error) {
	return GetRbacRolePermissionByPermissionIDWithPagination(repo.db, permissionID, limit, offset)
}

// GetRbacRolePermissionByRoleID runs a custom query, returning results as RbacRolePermission.
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByRoleID(roleID int64) ([]*RbacRolePermission, error) {
	return GetRbacRolePermissionByRoleID(repo.db, roleID)
}

// GetRbacRolePermissionByRoleIDCount runs a custom count query from repository
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByRoleIDCount(roleID int64) (int64, error) {
	return GetRbacRolePermissionByRoleIDCount(repo.db, roleID)
}

// GetRbacRolePermissionByRoleIDWithPagination runs a custom query with pagination from repository
func (repo *RbacRolePermissionRepository) GetRbacRolePermissionByRoleIDWithPagination(roleID int64, limit, offset int) ([]*RbacRolePermission, error) {
	return GetRbacRolePermissionByRoleIDWithPagination(repo.db, roleID, limit, offset)
}
