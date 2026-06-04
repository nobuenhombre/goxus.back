// Code generated from userstoken.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// IUsersTokenRepository defines the repository interface
type IUsersTokenRepository interface {
	Save(ut *UsersToken) error
	Delete(ut *UsersToken) error
	GetAll() ([]*UsersToken, error)
	GetLastID() (*UsersToken, error)
	GetUsersTokenByID(id int64) (*UsersToken, error)
	GetUsersTokenByToken(token string) (*UsersToken, error)
	FindAllUsersTokensByUserID(userID int64) ([]*UsersToken, error)
	DeleteExpiredTokens(ttlDays int) error
}

// Save saves the UsersToken to the database.
func (repo *UsersTokenRepository) Save(ut *UsersToken) error {
	return ut.Save(repo.db)
}

// Delete deletes the UsersToken from the database.
func (repo *UsersTokenRepository) Delete(ut *UsersToken) error {
	return ut.Delete(repo.db)
}

// UsersTokenRepository реализует работу с таблицей 'users_tokens'.
type UsersTokenRepository struct {
	db pgxdb.DBQuery
}

// NewUsersTokenRepository создает новый репозиторий.
func NewUsersTokenRepository(db pgxdb.DBQuery) *UsersTokenRepository {
	return &UsersTokenRepository{db: db}
}

// GetAll возвращает все записи
func (repo *UsersTokenRepository) GetAll() ([]*UsersToken, error) {
	return GetAllUsersToken(repo.db)
}

// GetLastID возвращает последний ID
func (repo *UsersTokenRepository) GetLastID() (*UsersToken, error) {
	return GetLastIDUsersToken(repo.db)
}

// GetUsersTokenByID возвращает одну запись по индексу 'users_tokens_pk'.
func (repo *UsersTokenRepository) GetUsersTokenByID(id int64) (*UsersToken, error) {
	return GetUsersTokenByID(repo.db, id)
}

// GetUsersTokenByToken возвращает одну запись по индексу 'users_tokens_token_uindex'.
func (repo *UsersTokenRepository) GetUsersTokenByToken(token string) (*UsersToken, error) {
	return GetUsersTokenByToken(repo.db, token)
}

// FindAllUsersTokensByUserID возвращает все записи по индексу 'users_tokens_user_id_index'.
func (repo *UsersTokenRepository) FindAllUsersTokensByUserID(userID int64) ([]*UsersToken, error) {
	return GetUsersTokensByUserID(repo.db, userID)
}

func (repo *UsersTokenRepository) DeleteExpiredTokens(ttlDays int) error {
	return DeleteExpiredTokens(repo.db, ttlDays)
}
