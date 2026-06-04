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
	GetAllWithPagination(limit, offset int) ([]*UsersToken, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*UsersToken, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*UsersToken, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*UsersToken, error)
	GetUsersTokenByID(id int64) (*UsersToken, error)
	GetUsersTokenByIDCount(id int64) (int64, error)
	GetUsersTokenByToken(token string) (*UsersToken, error)
	GetUsersTokenByTokenCount(token string) (int64, error)
	FindAllUsersTokensByUserID(userID int64) ([]*UsersToken, error)
	FindAllUsersTokensByUserIDWithPagination(userID int64, limit, offset int) ([]*UsersToken, error)
	FindAllUsersTokensByUserIDCount(userID int64) (int64, error)
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

// GetAllWithPagination возвращает записи с пагинацией
func (repo *UsersTokenRepository) GetAllWithPagination(limit, offset int) ([]*UsersToken, error) {
	return GetAllUsersTokenWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *UsersTokenRepository) GetAllCount() (int64, error) {
	return GetAllUsersTokenCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *UsersTokenRepository) GetBySQL(sqlstr string, args ...any) ([]*UsersToken, error) {
	return GetUsersTokensBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *UsersTokenRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*UsersToken, error) {
	return GetUsersTokensBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *UsersTokenRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetUsersTokensBySQLCount(repo.db, sqlstr, args...)
}

// GetLastID возвращает последний ID
func (repo *UsersTokenRepository) GetLastID() (*UsersToken, error) {
	return GetLastIDUsersToken(repo.db)
}

// GetUsersTokenByID возвращает одну запись по индексу 'users_tokens_pk'.
func (repo *UsersTokenRepository) GetUsersTokenByID(id int64) (*UsersToken, error) {
	return GetUsersTokenByID(repo.db, id)
}

// GetUsersTokenByIDCount возвращает количество записей по индексу 'users_tokens_pk'.
func (repo *UsersTokenRepository) GetUsersTokenByIDCount(id int64) (int64, error) {
	return GetUsersTokenByIDCount(repo.db, id)
}

// GetUsersTokenByToken возвращает одну запись по индексу 'users_tokens_token_uindex'.
func (repo *UsersTokenRepository) GetUsersTokenByToken(token string) (*UsersToken, error) {
	return GetUsersTokenByToken(repo.db, token)
}

// GetUsersTokenByTokenCount возвращает количество записей по индексу 'users_tokens_token_uindex'.
func (repo *UsersTokenRepository) GetUsersTokenByTokenCount(token string) (int64, error) {
	return GetUsersTokenByTokenCount(repo.db, token)
}

// FindAllUsersTokensByUserID возвращает все записи по индексу 'users_tokens_user_id_index'.
func (repo *UsersTokenRepository) FindAllUsersTokensByUserID(userID int64) ([]*UsersToken, error) {
	return GetUsersTokensByUserID(repo.db, userID)
}

// FindAllUsersTokensByUserIDWithPagination возвращает записи по индексу с пагинацией
func (repo *UsersTokenRepository) FindAllUsersTokensByUserIDWithPagination(userID int64, limit, offset int) ([]*UsersToken, error) {
	return GetUsersTokensByUserIDWithPagination(repo.db, userID, limit, offset)
}

// FindAllUsersTokensByUserIDCount возвращает количество записей по индексу 'users_tokens_user_id_index'.
func (repo *UsersTokenRepository) FindAllUsersTokensByUserIDCount(userID int64) (int64, error) {
	return GetUsersTokensByUserIDCount(repo.db, userID)
}

func (repo *UsersTokenRepository) DeleteExpiredTokens(ttlDays int) error {
	return DeleteExpiredTokens(repo.db, ttlDays)
}
