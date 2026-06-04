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
	GetLastID() (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	GetActiveUserByEmail(email string) (*User, error)
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

// GetLastID возвращает последний ID
func (repo *UserRepository) GetLastID() (*User, error) {
	return GetLastIDUser(repo.db)
}

// GetUserByEmail возвращает одну запись по индексу 'users_email_uindex'.
func (repo *UserRepository) GetUserByEmail(email string) (*User, error) {
	return GetUserByEmail(repo.db, email)
}

// GetUserByID возвращает одну запись по индексу 'users_pk'.
func (repo *UserRepository) GetUserByID(id int64) (*User, error) {
	return GetUserByID(repo.db, id)
}

// GetActiveUserByEmail runs a custom query, returning results as User.
func (repo *UserRepository) GetActiveUserByEmail(email string) (*User, error) {
	return GetActiveUserByEmail(repo.db, email)
}
