// Code generated from schemamigration.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// ISchemaMigrationRepository defines the repository interface
type ISchemaMigrationRepository interface {
	Save(sm *SchemaMigration) error
	Delete(sm *SchemaMigration) error
	GetAll() ([]*SchemaMigration, error)
	GetAllWithPagination(limit, offset int) ([]*SchemaMigration, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*SchemaMigration, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*SchemaMigration, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*SchemaMigration, error)
	GetSchemaMigrationByVersion(version int64) (*SchemaMigration, error)
	GetSchemaMigrationByVersionCount(version int64) (int64, error)
}

// Save saves the SchemaMigration to the database.
func (repo *SchemaMigrationRepository) Save(sm *SchemaMigration) error {
	return sm.Save(repo.db)
}

// Delete deletes the SchemaMigration from the database.
func (repo *SchemaMigrationRepository) Delete(sm *SchemaMigration) error {
	return sm.Delete(repo.db)
}

// SchemaMigrationRepository реализует работу с таблицей 'schema_migrations'.
type SchemaMigrationRepository struct {
	db pgxdb.DBQuery
}

// NewSchemaMigrationRepository создает новый репозиторий.
func NewSchemaMigrationRepository(db pgxdb.DBQuery) *SchemaMigrationRepository {
	return &SchemaMigrationRepository{db: db}
}

// GetAll возвращает все записи
func (repo *SchemaMigrationRepository) GetAll() ([]*SchemaMigration, error) {
	return GetAllSchemaMigration(repo.db)
}

// GetAllWithPagination возвращает записи с пагинацией
func (repo *SchemaMigrationRepository) GetAllWithPagination(limit, offset int) ([]*SchemaMigration, error) {
	return GetAllSchemaMigrationWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *SchemaMigrationRepository) GetAllCount() (int64, error) {
	return GetAllSchemaMigrationCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *SchemaMigrationRepository) GetBySQL(sqlstr string, args ...any) ([]*SchemaMigration, error) {
	return GetSchemaMigrationsBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *SchemaMigrationRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*SchemaMigration, error) {
	return GetSchemaMigrationsBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *SchemaMigrationRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetSchemaMigrationsBySQLCount(repo.db, sqlstr, args...)
}

// GetLastID возвращает последний ID
func (repo *SchemaMigrationRepository) GetLastID() (*SchemaMigration, error) {
	return GetLastIDSchemaMigration(repo.db)
}

// GetSchemaMigrationByVersion возвращает одну запись по индексу 'schema_migrations_pkey'.
func (repo *SchemaMigrationRepository) GetSchemaMigrationByVersion(version int64) (*SchemaMigration, error) {
	return GetSchemaMigrationByVersion(repo.db, version)
}

// GetSchemaMigrationByVersionCount возвращает количество записей по индексу 'schema_migrations_pkey'.
func (repo *SchemaMigrationRepository) GetSchemaMigrationByVersionCount(version int64) (int64, error) {
	return GetSchemaMigrationByVersionCount(repo.db, version)
}
