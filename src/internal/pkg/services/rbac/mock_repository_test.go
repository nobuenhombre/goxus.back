package rbac

import (
	"github.com/stretchr/testify/mock"
	"goxus/src/internal/pkg/db/goxus"
)

// mockRbacRoleRepo implements goxus.IRbacRoleRepository
type mockRbacRoleRepo struct {
	mock.Mock
}

func (m *mockRbacRoleRepo) Save(rr *goxus.RbacRole) error {
	args := m.Called(rr)
	return args.Error(0)
}

func (m *mockRbacRoleRepo) Delete(rr *goxus.RbacRole) error {
	args := m.Called(rr)
	return args.Error(0)
}

func (m *mockRbacRoleRepo) GetAll() ([]*goxus.RbacRole, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goxus.RbacRole), args.Error(1)
}

func (m *mockRbacRoleRepo) GetAllWithPagination(limit, offset int) ([]*goxus.RbacRole, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*goxus.RbacRole), args.Error(1)
}

func (m *mockRbacRoleRepo) GetAllCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacRoleRepo) GetBySQL(sqlstr string, args ...any) ([]*goxus.RbacRole, error) {
	callArgs := m.Called(append([]any{sqlstr}, args...)...)
	return callArgs.Get(0).([]*goxus.RbacRole), callArgs.Error(1)
}

func (m *mockRbacRoleRepo) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*goxus.RbacRole, error) {
	callArgs := m.Called(append([]any{sqlstr, limit, offset}, args...)...)
	return callArgs.Get(0).([]*goxus.RbacRole), callArgs.Error(1)
}

func (m *mockRbacRoleRepo) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	callArgs := m.Called(append([]any{sqlstr}, args...)...)
	return callArgs.Get(0).(int64), callArgs.Error(1)
}

func (m *mockRbacRoleRepo) GetLastID() (*goxus.RbacRole, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacRole), args.Error(1)
}

func (m *mockRbacRoleRepo) GetRbacRoleByID(id int64) (*goxus.RbacRole, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacRole), args.Error(1)
}

func (m *mockRbacRoleRepo) GetRbacRoleByIDCount(id int64) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacRoleRepo) GetRbacRoleBySlug(slug string) (*goxus.RbacRole, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacRole), args.Error(1)
}

func (m *mockRbacRoleRepo) GetRbacRoleBySlugCount(slug string) (int64, error) {
	args := m.Called(slug)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacRoleRepo) GetRolesByUserID(userID int64) ([]*goxus.RbacRole, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goxus.RbacRole), args.Error(1)
}

func (m *mockRbacRoleRepo) GetRolesByUserIDCount(userID int64) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacRoleRepo) GetRolesByUserIDWithPagination(userID int64, limit, offset int) ([]*goxus.RbacRole, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]*goxus.RbacRole), args.Error(1)
}

// mockRbacPermissionRepo implements goxus.IRbacPermissionRepository
type mockRbacPermissionRepo struct {
	mock.Mock
}

func (m *mockRbacPermissionRepo) Save(rp *goxus.RbacPermission) error {
	args := m.Called(rp)
	return args.Error(0)
}

func (m *mockRbacPermissionRepo) Delete(rp *goxus.RbacPermission) error {
	args := m.Called(rp)
	return args.Error(0)
}

func (m *mockRbacPermissionRepo) GetAll() ([]*goxus.RbacPermission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goxus.RbacPermission), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetAllWithPagination(limit, offset int) ([]*goxus.RbacPermission, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*goxus.RbacPermission), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetAllCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetBySQL(sqlstr string, args ...any) ([]*goxus.RbacPermission, error) {
	callArgs := m.Called(append([]any{sqlstr}, args...)...)
	return callArgs.Get(0).([]*goxus.RbacPermission), callArgs.Error(1)
}

func (m *mockRbacPermissionRepo) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*goxus.RbacPermission, error) {
	callArgs := m.Called(append([]any{sqlstr, limit, offset}, args...)...)
	return callArgs.Get(0).([]*goxus.RbacPermission), callArgs.Error(1)
}

func (m *mockRbacPermissionRepo) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	callArgs := m.Called(append([]any{sqlstr}, args...)...)
	return callArgs.Get(0).(int64), callArgs.Error(1)
}

func (m *mockRbacPermissionRepo) GetLastID() (*goxus.RbacPermission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacPermission), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetRbacPermissionByID(id int64) (*goxus.RbacPermission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacPermission), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetRbacPermissionByIDCount(id int64) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetRbacPermissionBySlug(slug string) (*goxus.RbacPermission, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacPermission), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetRbacPermissionBySlugCount(slug string) (int64, error) {
	args := m.Called(slug)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetPermissionsByRoleSlug(roleSlug string) ([]*goxus.RbacPermission, error) {
	args := m.Called(roleSlug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goxus.RbacPermission), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetPermissionsByRoleSlugCount(roleSlug string) (int64, error) {
	args := m.Called(roleSlug)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetPermissionsByRoleSlugWithPagination(roleSlug string, limit, offset int) ([]*goxus.RbacPermission, error) {
	args := m.Called(roleSlug, limit, offset)
	return args.Get(0).([]*goxus.RbacPermission), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetPermissionsByUserIDAndSlug(userID int64, permSlug string) ([]*goxus.RbacPermission, error) {
	args := m.Called(userID, permSlug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goxus.RbacPermission), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetPermissionsByUserIDAndSlugCount(userID int64, permSlug string) (int64, error) {
	args := m.Called(userID, permSlug)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacPermissionRepo) GetPermissionsByUserIDAndSlugWithPagination(userID int64, permSlug string, limit, offset int) ([]*goxus.RbacPermission, error) {
	args := m.Called(userID, permSlug, limit, offset)
	return args.Get(0).([]*goxus.RbacPermission), args.Error(1)
}

// mockRbacRolePermissionRepo implements goxus.IRbacRolePermissionRepository
type mockRbacRolePermissionRepo struct {
	mock.Mock
}

func (m *mockRbacRolePermissionRepo) Save(rrp *goxus.RbacRolePermission) error {
	args := m.Called(rrp)
	return args.Error(0)
}

func (m *mockRbacRolePermissionRepo) Delete(rrp *goxus.RbacRolePermission) error {
	args := m.Called(rrp)
	return args.Error(0)
}

func (m *mockRbacRolePermissionRepo) GetAll() ([]*goxus.RbacRolePermission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goxus.RbacRolePermission), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetAllWithPagination(limit, offset int) ([]*goxus.RbacRolePermission, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*goxus.RbacRolePermission), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetAllCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetBySQL(sqlstr string, args ...any) ([]*goxus.RbacRolePermission, error) {
	callArgs := m.Called(append([]any{sqlstr}, args...)...)
	return callArgs.Get(0).([]*goxus.RbacRolePermission), callArgs.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*goxus.RbacRolePermission, error) {
	callArgs := m.Called(append([]any{sqlstr, limit, offset}, args...)...)
	return callArgs.Get(0).([]*goxus.RbacRolePermission), callArgs.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	callArgs := m.Called(append([]any{sqlstr}, args...)...)
	return callArgs.Get(0).(int64), callArgs.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetLastID() (*goxus.RbacRolePermission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacRolePermission), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetRbacRolePermissionByID(id int64) (*goxus.RbacRolePermission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacRolePermission), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetRbacRolePermissionByIDCount(id int64) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetRbacRolePermissionByRoleIDPermissionID(roleID int64, permissionID int64) (*goxus.RbacRolePermission, error) {
	args := m.Called(roleID, permissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacRolePermission), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetRbacRolePermissionByRoleIDPermissionIDCount(roleID int64, permissionID int64) (int64, error) {
	args := m.Called(roleID, permissionID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetRbacRolePermissionByPermissionID(permissionID int64) ([]*goxus.RbacRolePermission, error) {
	args := m.Called(permissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goxus.RbacRolePermission), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetRbacRolePermissionByPermissionIDCount(permissionID int64) (int64, error) {
	args := m.Called(permissionID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetRbacRolePermissionByPermissionIDWithPagination(permissionID int64, limit, offset int) ([]*goxus.RbacRolePermission, error) {
	args := m.Called(permissionID, limit, offset)
	return args.Get(0).([]*goxus.RbacRolePermission), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetRbacRolePermissionByRoleID(roleID int64) ([]*goxus.RbacRolePermission, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goxus.RbacRolePermission), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetRbacRolePermissionByRoleIDCount(roleID int64) (int64, error) {
	args := m.Called(roleID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacRolePermissionRepo) GetRbacRolePermissionByRoleIDWithPagination(roleID int64, limit, offset int) ([]*goxus.RbacRolePermission, error) {
	args := m.Called(roleID, limit, offset)
	return args.Get(0).([]*goxus.RbacRolePermission), args.Error(1)
}

// mockRbacUserRoleRepo implements goxus.IRbacUserRoleRepository
type mockRbacUserRoleRepo struct {
	mock.Mock
}

func (m *mockRbacUserRoleRepo) Save(rur *goxus.RbacUserRole) error {
	args := m.Called(rur)
	return args.Error(0)
}

func (m *mockRbacUserRoleRepo) Delete(rur *goxus.RbacUserRole) error {
	args := m.Called(rur)
	return args.Error(0)
}

func (m *mockRbacUserRoleRepo) GetAll() ([]*goxus.RbacUserRole, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goxus.RbacUserRole), args.Error(1)
}

func (m *mockRbacUserRoleRepo) GetAllWithPagination(limit, offset int) ([]*goxus.RbacUserRole, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*goxus.RbacUserRole), args.Error(1)
}

func (m *mockRbacUserRoleRepo) GetAllCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacUserRoleRepo) GetBySQL(sqlstr string, args ...any) ([]*goxus.RbacUserRole, error) {
	callArgs := m.Called(append([]any{sqlstr}, args...)...)
	return callArgs.Get(0).([]*goxus.RbacUserRole), callArgs.Error(1)
}

func (m *mockRbacUserRoleRepo) GetBySQLWithPagination(sqlstr string, limit, offset int, args ...any) ([]*goxus.RbacUserRole, error) {
	callArgs := m.Called(append([]any{sqlstr, limit, offset}, args...)...)
	return callArgs.Get(0).([]*goxus.RbacUserRole), callArgs.Error(1)
}

func (m *mockRbacUserRoleRepo) GetBySQLCount(sqlstr string, args ...any) (int64, error) {
	callArgs := m.Called(append([]any{sqlstr}, args...)...)
	return callArgs.Get(0).(int64), callArgs.Error(1)
}

func (m *mockRbacUserRoleRepo) GetLastID() (*goxus.RbacUserRole, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacUserRole), args.Error(1)
}

func (m *mockRbacUserRoleRepo) GetRbacUserRoleByID(id int64) (*goxus.RbacUserRole, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacUserRole), args.Error(1)
}

func (m *mockRbacUserRoleRepo) GetRbacUserRoleByIDCount(id int64) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacUserRoleRepo) GetRbacUserRoleByUserIDRoleID(userID int64, roleID int64) (*goxus.RbacUserRole, error) {
	args := m.Called(userID, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goxus.RbacUserRole), args.Error(1)
}

func (m *mockRbacUserRoleRepo) GetRbacUserRoleByUserIDRoleIDCount(userID int64, roleID int64) (int64, error) {
	args := m.Called(userID, roleID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacUserRoleRepo) GetRbacUserRoleByRoleID(roleID int64) ([]*goxus.RbacUserRole, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*goxus.RbacUserRole), args.Error(1)
}

func (m *mockRbacUserRoleRepo) GetRbacUserRoleByRoleIDCount(roleID int64) (int64, error) {
	args := m.Called(roleID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRbacUserRoleRepo) GetRbacUserRoleByRoleIDWithPagination(roleID int64, limit, offset int) ([]*goxus.RbacUserRole, error) {
	args := m.Called(roleID, limit, offset)
	return args.Get(0).([]*goxus.RbacUserRole), args.Error(1)
}
