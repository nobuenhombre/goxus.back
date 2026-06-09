// Code generated from settingstype.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// ISettingsTypeRepository defines the repository interface
type ISettingsTypeRepository interface {
	Save(st *SettingsType) error
	Delete(st *SettingsType) error
	GetAll() ([]*SettingsType, error)
	GetAllWithPagination(limit, offset int) ([]*SettingsType, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*SettingsType, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*SettingsType, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*SettingsType, error)
	GetSettingsTypeByName(name string) (*SettingsType, error)
	GetSettingsTypeByID(id int64) (*SettingsType, error)
}

// Save saves the SettingsType to the database.
func (repo *SettingsTypeRepository) Save(st *SettingsType) error {
	return st.Save(repo.db)
}

// Delete deletes the SettingsType from the database.
func (repo *SettingsTypeRepository) Delete(st *SettingsType) error {
	return st.Delete(repo.db)
}

// SettingsTypeRepository реализует работу с таблицей 'settings_types'.
type SettingsTypeRepository struct {
	db pgxdb.DBQuery
}

// NewSettingsTypeRepository создает новый репозиторий.
func NewSettingsTypeRepository(db pgxdb.DBQuery) *SettingsTypeRepository {
	return &SettingsTypeRepository{db: db}
}

// GetAll возвращает все записи
func (repo *SettingsTypeRepository) GetAll() ([]*SettingsType, error) {
	return GetAllSettingsType(repo.db)
}

// GetAllWithPagination возвращает записи с пагинацией
func (repo *SettingsTypeRepository) GetAllWithPagination(limit, offset int) ([]*SettingsType, error) {
	return GetAllSettingsTypeWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *SettingsTypeRepository) GetAllCount() (int64, error) {
	return GetAllSettingsTypeCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *SettingsTypeRepository) GetBySQL(sqlstr string, args ...any) ([]*SettingsType, error) {
	return GetSettingsTypesBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *SettingsTypeRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*SettingsType, error) {
	return GetSettingsTypesBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *SettingsTypeRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetSettingsTypesBySQLCount(repo.db, sqlstr, args...)
}

// GetLastID возвращает последний ID
func (repo *SettingsTypeRepository) GetLastID() (*SettingsType, error) {
	return GetLastIDSettingsType(repo.db)
}

// GetSettingsTypeByName возвращает одну запись по индексу 'settings_types_name_uindex'.
func (repo *SettingsTypeRepository) GetSettingsTypeByName(name string) (*SettingsType, error) {
	return GetSettingsTypeByName(repo.db, name)
}

// GetSettingsTypeByID возвращает одну запись по индексу 'settings_types_pk'.
func (repo *SettingsTypeRepository) GetSettingsTypeByID(id int64) (*SettingsType, error) {
	return GetSettingsTypeByID(repo.db, id)
}
