// Code generated from user.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// IUserRepository defines the repository interface
type IUserRepository interface {
	Save(u *User) error
	Delete(u *User) error
	GetAll() ([]*User, error)
	GetAllWithPagination(limit, offset int) ([]*User, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*User, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*User, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByEmailCount(email string) (int64, error)
	GetUserByID(id int64) (*User, error)
	GetUserByIDCount(id int64) (int64, error)
	GetActiveUserByEmail(email string) (*User, error)
	GetActiveUserByEmailCount(email string) (int64, error)
}

// Save saves the User to the database.
func (repo *UserRepository) Save(u *User) error {
	return u.Save(repo.db)
}

// Delete deletes the User from the database.
func (repo *UserRepository) Delete(u *User) error {
	return u.Delete(repo.db)
}

// UserRepository реализует работу с таблицей 'users'.
type UserRepository struct {
	db pgxdb.DBQuery
}

// NewUserRepository создает новый репозиторий.
func NewUserRepository(db pgxdb.DBQuery) *UserRepository {
	return &UserRepository{db: db}
}

// GetAll возвращает все записи
func (repo *UserRepository) GetAll() ([]*User, error) {
	return GetAllUser(repo.db)
}

// GetAllWithPagination возвращает записи с пагинацией
func (repo *UserRepository) GetAllWithPagination(limit, offset int) ([]*User, error) {
	return GetAllUserWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *UserRepository) GetAllCount() (int64, error) {
	return GetAllUserCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *UserRepository) GetBySQL(sqlstr string, args ...any) ([]*User, error) {
	return GetUsersBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *UserRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*User, error) {
	return GetUsersBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *UserRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetUsersBySQLCount(repo.db, sqlstr, args...)
}

// GetLastID возвращает последний ID
func (repo *UserRepository) GetLastID() (*User, error) {
	return GetLastIDUser(repo.db)
}

// GetUserByEmail возвращает одну запись по индексу 'users_email_uindex'.
func (repo *UserRepository) GetUserByEmail(email string) (*User, error) {
	return GetUserByEmail(repo.db, email)
}

// GetUserByEmailCount возвращает количество записей по индексу 'users_email_uindex'.
func (repo *UserRepository) GetUserByEmailCount(email string) (int64, error) {
	return GetUserByEmailCount(repo.db, email)
}

// GetUserByID возвращает одну запись по индексу 'users_pk'.
func (repo *UserRepository) GetUserByID(id int64) (*User, error) {
	return GetUserByID(repo.db, id)
}

// GetUserByIDCount возвращает количество записей по индексу 'users_pk'.
func (repo *UserRepository) GetUserByIDCount(id int64) (int64, error) {
	return GetUserByIDCount(repo.db, id)
}

// GetActiveUserByEmail runs a custom query, returning results as User.
func (repo *UserRepository) GetActiveUserByEmail(email string) (*User, error) {
	return GetActiveUserByEmail(repo.db, email)
}

// GetActiveUserByEmailCount runs a custom count query from repository
func (repo *UserRepository) GetActiveUserByEmailCount(email string) (int64, error) {
	return GetActiveUserByEmailCount(repo.db, email)
}
