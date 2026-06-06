// Code generated from userwithrole.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// IUserWithRoleRepository defines the repository interface
type IUserWithRoleRepository interface {
	GetAll() ([]*UserWithRole, error)
	GetAllWithPagination(limit, offset int) ([]*UserWithRole, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*UserWithRole, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*UserWithRole, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*UserWithRole, error)
}

// UserWithRoleRepository реализует работу с таблицей 'user_with_roles'.
type UserWithRoleRepository struct {
	db pgxdb.DBQuery
}

// NewUserWithRoleRepository создает новый репозиторий.
func NewUserWithRoleRepository(db pgxdb.DBQuery) *UserWithRoleRepository {
	return &UserWithRoleRepository{db: db}
}

// GetAll возвращает все записи
func (repo *UserWithRoleRepository) GetAll() ([]*UserWithRole, error) {
	return GetAllUserWithRole(repo.db)
}

// GetAllWithPagination возвращает записи с пагинацией
func (repo *UserWithRoleRepository) GetAllWithPagination(limit, offset int) ([]*UserWithRole, error) {
	return GetAllUserWithRoleWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *UserWithRoleRepository) GetAllCount() (int64, error) {
	return GetAllUserWithRoleCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *UserWithRoleRepository) GetBySQL(sqlstr string, args ...any) ([]*UserWithRole, error) {
	return GetUserWithRolesBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *UserWithRoleRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*UserWithRole, error) {
	return GetUserWithRolesBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *UserWithRoleRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetUserWithRolesBySQLCount(repo.db, sqlstr, args...)
}

// GetLastID возвращает последний ID
func (repo *UserWithRoleRepository) GetLastID() (*UserWithRole, error) {
	return GetLastIDUserWithRole(repo.db)
}
