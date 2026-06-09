// Code generated from userssetting.xo-xouid.go. DO NOT EDIT.
package goxus

import (
	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
)

// IUsersSettingRepository defines the repository interface
type IUsersSettingRepository interface {
	Save(us *UsersSetting) error
	Delete(us *UsersSetting) error
	GetAll() ([]*UsersSetting, error)
	GetAllWithPagination(limit, offset int) ([]*UsersSetting, error)
	GetAllCount() (int64, error)
	GetBySQL(sqlstr string, args ...any) ([]*UsersSetting, error)
	GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*UsersSetting, error)
	GetBySQLCount(sqlstr string, args ...any) (int64, error)
	GetLastID() (*UsersSetting, error)
	GetUsersSettingByID(id int64) (*UsersSetting, error)
	GetUsersSettingByUserIDSettingsID(userID int64, settingsID int64) (*UsersSetting, error)
}

// Save saves the UsersSetting to the database.
func (repo *UsersSettingRepository) Save(us *UsersSetting) error {
	return us.Save(repo.db)
}

// Delete deletes the UsersSetting from the database.
func (repo *UsersSettingRepository) Delete(us *UsersSetting) error {
	return us.Delete(repo.db)
}

// UsersSettingRepository реализует работу с таблицей 'users_settings'.
type UsersSettingRepository struct {
	db pgxdb.DBQuery
}

// NewUsersSettingRepository создает новый репозиторий.
func NewUsersSettingRepository(db pgxdb.DBQuery) *UsersSettingRepository {
	return &UsersSettingRepository{db: db}
}

// GetAll возвращает все записи
func (repo *UsersSettingRepository) GetAll() ([]*UsersSetting, error) {
	return GetAllUsersSetting(repo.db)
}

// GetAllWithPagination возвращает записи с пагинацией
func (repo *UsersSettingRepository) GetAllWithPagination(limit, offset int) ([]*UsersSetting, error) {
	return GetAllUsersSettingWithPagination(repo.db, limit, offset)
}

// GetAllCount возвращает количество записей
func (repo *UsersSettingRepository) GetAllCount() (int64, error) {
	return GetAllUsersSettingCount(repo.db)
}

// GetBySQL возвращает записи по произвольному SQL
func (repo *UsersSettingRepository) GetBySQL(sqlstr string, args ...any) ([]*UsersSetting, error) {
	return GetUsersSettingsBySQL(repo.db, sqlstr, args...)
}

// GetBySQLWithPagination возвращает записи по произвольному SQL с пагинацией
func (repo *UsersSettingRepository) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*UsersSetting, error) {
	return GetUsersSettingsBySQLWithPagination(repo.db, sqlstr, limit, offset, args...)
}

// GetBySQLCount возвращает количество записей по произвольному SQL
func (repo *UsersSettingRepository) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	return GetUsersSettingsBySQLCount(repo.db, sqlstr, args...)
}

// GetLastID возвращает последний ID
func (repo *UsersSettingRepository) GetLastID() (*UsersSetting, error) {
	return GetLastIDUsersSetting(repo.db)
}

// GetUsersSettingByID возвращает одну запись по индексу 'users_settings_pk'.
func (repo *UsersSettingRepository) GetUsersSettingByID(id int64) (*UsersSetting, error) {
	return GetUsersSettingByID(repo.db, id)
}

// GetUsersSettingByUserIDSettingsID возвращает одну запись по индексу 'users_settings_user_id_settings_id_idx'.
func (repo *UsersSettingRepository) GetUsersSettingByUserIDSettingsID(userID int64, settingsID int64) (*UsersSetting, error) {
	return GetUsersSettingByUserIDSettingsID(repo.db, userID, settingsID)
}
