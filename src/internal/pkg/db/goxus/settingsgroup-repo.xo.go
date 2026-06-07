// Code generated from settingsgroup.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// ISettingsGroupRepository defines the repository interface
type ISettingsGroupRepository interface {
	Save(sg *SettingsGroup) error
	Delete(sg *SettingsGroup) error
	GetAll() ([]*SettingsGroup, error)
	GetAllWithPagination(limit, offset int) ([]*SettingsGroup, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*SettingsGroup, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*SettingsGroup, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*SettingsGroup, error)
	GetSettingsGroupByID(id int64) (*SettingsGroup, error)
	GetSettingsGroupByIDCount(id int64) (int64, error)
}

// Save saves the SettingsGroup to the database.
func (repo *SettingsGroupRepository) Save(sg *SettingsGroup) error {
	return sg.Save(repo.db)
}

// Delete deletes the SettingsGroup from the database.
func (repo *SettingsGroupRepository) Delete(sg *SettingsGroup) error {
	return sg.Delete(repo.db)
}

// SettingsGroupRepository реализует работу с таблицей 'settings_groups'.
type SettingsGroupRepository struct {
	db pgxdb.DBQuery
}

// NewSettingsGroupRepository создает новый репозиторий.
func NewSettingsGroupRepository(db pgxdb.DBQuery) *SettingsGroupRepository {
	return &SettingsGroupRepository{db: db}
}

// GetAll возвращает все записи
func (repo *SettingsGroupRepository) GetAll() ([]*SettingsGroup, error) {
	return GetAllSettingsGroup(repo.db)
}

// GetAllWithPagination возвращает записи с пагинацией
func (repo *SettingsGroupRepository) GetAllWithPagination(limit, offset int) ([]*SettingsGroup, error) {
	return GetAllSettingsGroupWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *SettingsGroupRepository) GetAllCount() (int64, error) {
	return GetAllSettingsGroupCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *SettingsGroupRepository) GetBySQL(sqlstr string, args ...any) ([]*SettingsGroup, error) {
	return GetSettingsGroupsBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *SettingsGroupRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*SettingsGroup, error) {
	return GetSettingsGroupsBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *SettingsGroupRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetSettingsGroupsBySQLCount(repo.db, sqlstr, args...)
}

// GetLastID возвращает последний ID
func (repo *SettingsGroupRepository) GetLastID() (*SettingsGroup, error) {
	return GetLastIDSettingsGroup(repo.db)
}

// GetSettingsGroupByID возвращает одну запись по индексу 'settings_groups_pk'.
func (repo *SettingsGroupRepository) GetSettingsGroupByID(id int64) (*SettingsGroup, error) {
	return GetSettingsGroupByID(repo.db, id)
}

// GetSettingsGroupByIDCount возвращает количество записей по индексу 'settings_groups_pk'.
func (repo *SettingsGroupRepository) GetSettingsGroupByIDCount(id int64) (int64, error) {
	return GetSettingsGroupByIDCount(repo.db, id)
}
