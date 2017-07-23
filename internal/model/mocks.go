package model

import (
	"context"
	"database/sql"
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/stretchr/testify/mock"
)

type MockGroupProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: _a0, _a1
func (_m *MockGroupProvider) Insert(_a0 context.Context, _a1 *GroupEntity) (*GroupEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, *GroupEntity) *GroupEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *GroupEntity) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUserID provides a mock function with given fields: _a0, _a1
func (_m *MockGroupProvider) FindByUserID(_a0 context.Context, _a1 int64) ([]*GroupEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*GroupEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByID provides a mock function with given fields: _a0, _a1
func (_m *MockGroupProvider) FindOneByID(_a0 context.Context, _a1 int64) (*GroupEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64) *GroupEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: _a0, _a1
func (_m *MockGroupProvider) Find(_a0 context.Context, _a1 *GroupFindExpr) ([]*GroupEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, *GroupFindExpr) []*GroupEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *GroupFindExpr) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, createdBy, name, description
func (_m *MockGroupProvider) Create(ctx context.Context, createdBy int64, name string, description *ntypes.String) (*GroupEntity, error) {
	ret := _m.Called(ctx, createdBy, name, description)

	var r0 *GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, *ntypes.String) *GroupEntity); ok {
		r0 = rf(ctx, createdBy, name, description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, string, *ntypes.String) error); ok {
		r1 = rf(ctx, createdBy, name, description)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOneByID provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockGroupProvider) UpdateOneByID(_a0 context.Context, _a1 int64, _a2 *GroupPatch) (*GroupEntity, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *GroupEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64, *GroupPatch) *GroupEntity); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, *GroupPatch) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOneByID provides a mock function with given fields: _a0, _a1
func (_m *MockGroupProvider) DeleteOneByID(_a0 context.Context, _a1 int64) (int64, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsGranted provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockGroupProvider) IsGranted(_a0 context.Context, _a1 int64, _a2 charon.Permission) (bool, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64, charon.Permission) bool); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, charon.Permission) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPermissions provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockGroupProvider) SetPermissions(_a0 context.Context, _a1 int64, _a2 ...charon.Permission) (int64, int64, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64, ...charon.Permission) int64); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, int64, ...charon.Permission) int64); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, int64, ...charon.Permission) error); ok {
		r2 = rf(_a0, _a1, _a2...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type MockGroupPermissionsProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: _a0, _a1
func (_m *MockGroupPermissionsProvider) Insert(_a0 context.Context, _a1 *GroupPermissionsEntity) (*GroupPermissionsEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *GroupPermissionsEntity
	if rf, ok := ret.Get(0).(func(context.Context, *GroupPermissionsEntity) *GroupPermissionsEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupPermissionsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *GroupPermissionsEntity) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type MockPermissionProvider struct {
	mock.Mock
}

// Find provides a mock function with given fields: ctx, criteria
func (_m *MockPermissionProvider) Find(ctx context.Context, criteria *PermissionFindExpr) ([]*PermissionEntity, error) {
	ret := _m.Called(ctx, criteria)

	var r0 []*PermissionEntity
	if rf, ok := ret.Get(0).(func(context.Context, *PermissionFindExpr) []*PermissionEntity); ok {
		r0 = rf(ctx, criteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*PermissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *PermissionFindExpr) error); ok {
		r1 = rf(ctx, criteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByID provides a mock function with given fields: ctx, id
func (_m *MockPermissionProvider) FindOneByID(ctx context.Context, id int64) (*PermissionEntity, error) {
	ret := _m.Called(ctx, id)

	var r0 *PermissionEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64) *PermissionEntity); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PermissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUserID provides a mock function with given fields: ctx, userID
func (_m *MockPermissionProvider) FindByUserID(ctx context.Context, userID int64) ([]*PermissionEntity, error) {
	ret := _m.Called(ctx, userID)

	var r0 []*PermissionEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*PermissionEntity); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*PermissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByGroupID provides a mock function with given fields: ctx, groupID
func (_m *MockPermissionProvider) FindByGroupID(ctx context.Context, groupID int64) ([]*PermissionEntity, error) {
	ret := _m.Called(ctx, groupID)

	var r0 []*PermissionEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*PermissionEntity); ok {
		r0 = rf(ctx, groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*PermissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: ctx, permissions
func (_m *MockPermissionProvider) Register(ctx context.Context, permissions charon.Permissions) (int64, int64, int64, error) {
	ret := _m.Called(ctx, permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, charon.Permissions) int64); ok {
		r0 = rf(ctx, permissions)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, charon.Permissions) int64); ok {
		r1 = rf(ctx, permissions)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 int64
	if rf, ok := ret.Get(2).(func(context.Context, charon.Permissions) int64); ok {
		r2 = rf(ctx, permissions)
	} else {
		r2 = ret.Get(2).(int64)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, charon.Permissions) error); ok {
		r3 = rf(ctx, permissions)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// Insert provides a mock function with given fields: ctx, entity
func (_m *MockPermissionProvider) Insert(ctx context.Context, entity *PermissionEntity) (*PermissionEntity, error) {
	ret := _m.Called(ctx, entity)

	var r0 *PermissionEntity
	if rf, ok := ret.Get(0).(func(context.Context, *PermissionEntity) *PermissionEntity); ok {
		r0 = rf(ctx, entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PermissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *PermissionEntity) error); ok {
		r1 = rf(ctx, entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type MockPermissionRegistry struct {
	mock.Mock
}

// Exists provides a mock function with given fields: ctx, permission
func (_m *MockPermissionRegistry) Exists(ctx context.Context, permission charon.Permission) bool {
	ret := _m.Called(ctx, permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, charon.Permission) bool); ok {
		r0 = rf(ctx, permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Register provides a mock function with given fields: ctx, permissions
func (_m *MockPermissionRegistry) Register(ctx context.Context, permissions charon.Permissions) (int64, int64, int64, error) {
	ret := _m.Called(ctx, permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, charon.Permissions) int64); ok {
		r0 = rf(ctx, permissions)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, charon.Permissions) int64); ok {
		r1 = rf(ctx, permissions)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 int64
	if rf, ok := ret.Get(2).(func(context.Context, charon.Permissions) int64); ok {
		r2 = rf(ctx, permissions)
	} else {
		r2 = ret.Get(2).(int64)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, charon.Permissions) error); ok {
		r3 = rf(ctx, permissions)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

type MockRows struct {
	mock.Mock
}

// ColumnTypes provides a mock function with given fields:
func (_m *MockRows) ColumnTypes() ([]*sql.ColumnType, error) {
	ret := _m.Called()

	var r0 []*sql.ColumnType
	if rf, ok := ret.Get(0).(func() []*sql.ColumnType); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*sql.ColumnType)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Columns provides a mock function with given fields:
func (_m *MockRows) Columns() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Err provides a mock function with given fields:
func (_m *MockRows) Err() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Next provides a mock function with given fields:
func (_m *MockRows) Next() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NextResultSet provides a mock function with given fields:
func (_m *MockRows) NextResultSet() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Scan provides a mock function with given fields: dest
func (_m *MockRows) Scan(dest ...interface{}) error {
	ret := _m.Called(dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(...interface{}) error); ok {
		r0 = rf(dest...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type MockCompositionWriter struct {
	mock.Mock
}

// WriteComposition provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockCompositionWriter) WriteComposition(_a0 string, _a1 *Composer, _a2 *CompositionOpts) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *Composer, *CompositionOpts) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockSuite struct {
	mock.Mock
}

// setup provides a mock function with given fields: _a0
func (_m *mockSuite) setup(_a0 testing.T) {
	_m.Called(_a0)
}

// teardown provides a mock function with given fields: _a0
func (_m *mockSuite) teardown(_a0 testing.T) {
	_m.Called(_a0)
}

type MockUserProvider struct {
	mock.Mock
}

// Exists provides a mock function with given fields: _a0, _a1
func (_m *MockUserProvider) Exists(_a0 context.Context, _a1 int64) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *MockUserProvider) Create(_a0 context.Context, _a1 *UserEntity) (*UserEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(context.Context, *UserEntity) *UserEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *UserEntity) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: _a0, _a1
func (_m *MockUserProvider) Insert(_a0 context.Context, _a1 *UserEntity) (*UserEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(context.Context, *UserEntity) *UserEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *UserEntity) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSuperuser provides a mock function with given fields: ctx, username, password, FirstName, LastName
func (_m *MockUserProvider) CreateSuperuser(ctx context.Context, username string, password []byte, FirstName string, LastName string) (*UserEntity, error) {
	ret := _m.Called(ctx, username, password, FirstName, LastName)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte, string, string) *UserEntity); ok {
		r0 = rf(ctx, username, password, FirstName, LastName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []byte, string, string) error); ok {
		r1 = rf(ctx, username, password, FirstName, LastName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Count provides a mock function with given fields: _a0
func (_m *MockUserProvider) Count(_a0 context.Context) (int64, error) {
	ret := _m.Called(_a0)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context) int64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateLastLoginAt provides a mock function with given fields: ctx, id
func (_m *MockUserProvider) UpdateLastLoginAt(ctx context.Context, id int64) (int64, error) {
	ret := _m.Called(ctx, id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChangePassword provides a mock function with given fields: ctx, id, password
func (_m *MockUserProvider) ChangePassword(ctx context.Context, id int64, password string) error {
	ret := _m.Called(ctx, id, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: _a0, _a1
func (_m *MockUserProvider) Find(_a0 context.Context, _a1 *UserFindExpr) ([]*UserEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*UserEntity
	if rf, ok := ret.Get(0).(func(context.Context, *UserFindExpr) []*UserEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *UserFindExpr) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByID provides a mock function with given fields: _a0, _a1
func (_m *MockUserProvider) FindOneByID(_a0 context.Context, _a1 int64) (*UserEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64) *UserEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByUsername provides a mock function with given fields: _a0, _a1
func (_m *MockUserProvider) FindOneByUsername(_a0 context.Context, _a1 string) (*UserEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(context.Context, string) *UserEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOneByID provides a mock function with given fields: _a0, _a1
func (_m *MockUserProvider) DeleteOneByID(_a0 context.Context, _a1 int64) (int64, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOneByID provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockUserProvider) UpdateOneByID(_a0 context.Context, _a1 int64, _a2 *UserPatch) (*UserEntity, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(context.Context, int64, *UserPatch) *UserEntity); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, *UserPatch) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegistrationConfirmation provides a mock function with given fields: ctx, id, confirmationToken
func (_m *MockUserProvider) RegistrationConfirmation(ctx context.Context, id int64, confirmationToken string) (int64, error) {
	ret := _m.Called(ctx, id, confirmationToken)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) int64); ok {
		r0 = rf(ctx, id, confirmationToken)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, string) error); ok {
		r1 = rf(ctx, id, confirmationToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsGranted provides a mock function with given fields: ctx, id, permission
func (_m *MockUserProvider) IsGranted(ctx context.Context, id int64, permission charon.Permission) (bool, error) {
	ret := _m.Called(ctx, id, permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64, charon.Permission) bool); ok {
		r0 = rf(ctx, id, permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, charon.Permission) error); ok {
		r1 = rf(ctx, id, permission)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPermissions provides a mock function with given fields: ctx, id, permissions
func (_m *MockUserProvider) SetPermissions(ctx context.Context, id int64, permissions ...charon.Permission) (int64, int64, error) {
	ret := _m.Called(ctx, id, permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64, ...charon.Permission) int64); ok {
		r0 = rf(ctx, id, permissions...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, int64, ...charon.Permission) int64); ok {
		r1 = rf(ctx, id, permissions...)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, int64, ...charon.Permission) error); ok {
		r2 = rf(ctx, id, permissions...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type MockUserGroupsProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: ctx, ent
func (_m *MockUserGroupsProvider) Insert(ctx context.Context, ent *UserGroupsEntity) (*UserGroupsEntity, error) {
	ret := _m.Called(ctx, ent)

	var r0 *UserGroupsEntity
	if rf, ok := ret.Get(0).(func(context.Context, *UserGroupsEntity) *UserGroupsEntity); ok {
		r0 = rf(ctx, ent)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserGroupsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *UserGroupsEntity) error); ok {
		r1 = rf(ctx, ent)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: ctx, userID, groupID
func (_m *MockUserGroupsProvider) Exists(ctx context.Context, userID int64, groupID int64) (bool, error) {
	ret := _m.Called(ctx, userID, groupID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) bool); ok {
		r0 = rf(ctx, userID, groupID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, userID, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: ctx, expr
func (_m *MockUserGroupsProvider) Find(ctx context.Context, expr *UserGroupsFindExpr) ([]*UserGroupsEntity, error) {
	ret := _m.Called(ctx, expr)

	var r0 []*UserGroupsEntity
	if rf, ok := ret.Get(0).(func(context.Context, *UserGroupsFindExpr) []*UserGroupsEntity); ok {
		r0 = rf(ctx, expr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*UserGroupsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *UserGroupsFindExpr) error); ok {
		r1 = rf(ctx, expr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: ctx, userID, groupIDs
func (_m *MockUserGroupsProvider) Set(ctx context.Context, userID int64, groupIDs []int64) (int64, int64, error) {
	ret := _m.Called(ctx, userID, groupIDs)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64, []int64) int64); ok {
		r0 = rf(ctx, userID, groupIDs)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, int64, []int64) int64); ok {
		r1 = rf(ctx, userID, groupIDs)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, int64, []int64) error); ok {
		r2 = rf(ctx, userID, groupIDs)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteByUserID provides a mock function with given fields: ctx, id
func (_m *MockUserGroupsProvider) DeleteByUserID(ctx context.Context, id int64) (int64, error) {
	ret := _m.Called(ctx, id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type MockUserPermissionsProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: _a0, _a1
func (_m *MockUserPermissionsProvider) Insert(_a0 context.Context, _a1 *UserPermissionsEntity) (*UserPermissionsEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *UserPermissionsEntity
	if rf, ok := ret.Get(0).(func(context.Context, *UserPermissionsEntity) *UserPermissionsEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserPermissionsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *UserPermissionsEntity) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByUserID provides a mock function with given fields: _a0, _a1
func (_m *MockUserPermissionsProvider) DeleteByUserID(_a0 context.Context, _a1 int64) (int64, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
