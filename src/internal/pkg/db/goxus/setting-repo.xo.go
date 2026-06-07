// Code generated from setting.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// ISettingRepository defines the repository interface
type ISettingRepository interface {
	Save(s *Setting) error
	Delete(s *Setting) error
	GetAll() ([]*Setting, error)
	GetAllWithPagination(limit, offset int) ([]*Setting, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*Setting, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*Setting, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*Setting, error)
	GetSettingByID(id int64) (*Setting, error)
	GetSettingByIDCount(id int64) (int64, error)
}

// Save saves the Setting to the database.
func (repo *SettingRepository) Save(s *Setting) error {
	return s.Save(repo.db)
}

// Delete deletes the Setting from the database.
func (repo *SettingRepository) Delete(s *Setting) error {
	return s.Delete(repo.db)
}

// SettingRepository реализует работу с таблицей 'settings'.
type SettingRepository struct {
	db pgxdb.DBQuery
}

// NewSettingRepository создает новый репозиторий.
func NewSettingRepository(db pgxdb.DBQuery) *SettingRepository {
	return &SettingRepository{db: db}
}

// GetAll возвращает все записи
func (repo *SettingRepository) GetAll() ([]*Setting, error) {
	return GetAllSetting(repo.db)
}

// GetAllWithPagination возвращает записи с пагинацией
func (repo *SettingRepository) GetAllWithPagination(limit, offset int) ([]*Setting, error) {
	return GetAllSettingWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *SettingRepository) GetAllCount() (int64, error) {
	return GetAllSettingCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *SettingRepository) GetBySQL(sqlstr string, args ...any) ([]*Setting, error) {
	return GetSettingsBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *SettingRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*Setting, error) {
	return GetSettingsBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *SettingRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetSettingsBySQLCount(repo.db, sqlstr, args...)
}

// GetLastID возвращает последний ID
func (repo *SettingRepository) GetLastID() (*Setting, error) {
	return GetLastIDSetting(repo.db)
}

// GetSettingByID возвращает одну запись по индексу 'settings_pk'.
func (repo *SettingRepository) GetSettingByID(id int64) (*Setting, error) {
	return GetSettingByID(repo.db, id)
}

// GetSettingByIDCount возвращает количество записей по индексу 'settings_pk'.
func (repo *SettingRepository) GetSettingByIDCount(id int64) (int64, error) {
	return GetSettingByIDCount(repo.db, id)
}
