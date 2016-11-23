package model

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/stretchr/testify/mock"
)

import "github.com/piotrkowalczuk/ntypes"

type MockGroupProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: entity
func (_m *MockGroupProvider) Insert(entity *GroupEntity) (*GroupEntity, error) {
	ret := _m.Called(entity)

	var r0 *GroupEntity
	if rf, ok := ret.Get(0).(func(*GroupEntity) *GroupEntity); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*GroupEntity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUserID provides a mock function with given fields: _a0
func (_m *MockGroupProvider) FindByUserID(_a0 int64) ([]*GroupEntity, error) {
	ret := _m.Called(_a0)

	var r0 []*GroupEntity
	if rf, ok := ret.Get(0).(func(int64) []*GroupEntity); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByID provides a mock function with given fields: _a0
func (_m *MockGroupProvider) FindOneByID(_a0 int64) (*GroupEntity, error) {
	ret := _m.Called(_a0)

	var r0 *GroupEntity
	if rf, ok := ret.Get(0).(func(int64) *GroupEntity); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: c
func (_m *MockGroupProvider) Find(c *GroupCriteria) ([]*GroupEntity, error) {
	ret := _m.Called(c)

	var r0 []*GroupEntity
	if rf, ok := ret.Get(0).(func(*GroupCriteria) []*GroupEntity); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*GroupCriteria) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: createdBy, name, description
func (_m *MockGroupProvider) Create(createdBy int64, name string, description *ntypes.String) (*GroupEntity, error) {
	ret := _m.Called(createdBy, name, description)

	var r0 *GroupEntity
	if rf, ok := ret.Get(0).(func(int64, string, *ntypes.String) *GroupEntity); ok {
		r0 = rf(createdBy, name, description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, string, *ntypes.String) error); ok {
		r1 = rf(createdBy, name, description)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOneByID provides a mock function with given fields: id, updatedBy, name, description
func (_m *MockGroupProvider) UpdateOneByID(id int64, updatedBy int64, name *ntypes.String, description *ntypes.String) (*GroupEntity, error) {
	ret := _m.Called(id, updatedBy, name, description)

	var r0 *GroupEntity
	if rf, ok := ret.Get(0).(func(int64, int64, *ntypes.String, *ntypes.String) *GroupEntity); ok {
		r0 = rf(id, updatedBy, name, description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64, *ntypes.String, *ntypes.String) error); ok {
		r1 = rf(id, updatedBy, name, description)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOneByID provides a mock function with given fields: id
func (_m *MockGroupProvider) DeleteOneByID(id int64) (int64, error) {
	ret := _m.Called(id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64) int64); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsGranted provides a mock function with given fields: id, permission
func (_m *MockGroupProvider) IsGranted(id int64, permission charon.Permission) (bool, error) {
	ret := _m.Called(id, permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int64, charon.Permission) bool); ok {
		r0 = rf(id, permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, charon.Permission) error); ok {
		r1 = rf(id, permission)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPermissions provides a mock function with given fields: id, permissions
func (_m *MockGroupProvider) SetPermissions(id int64, permissions ...charon.Permission) (int64, int64, error) {
	ret := _m.Called(id, permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, ...charon.Permission) int64); ok {
		r0 = rf(id, permissions...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(int64, ...charon.Permission) int64); ok {
		r1 = rf(id, permissions...)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int64, ...charon.Permission) error); ok {
		r2 = rf(id, permissions...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type MockGroupPermissionsProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: entity
func (_m *MockGroupPermissionsProvider) Insert(entity *GroupPermissionsEntity) (*GroupPermissionsEntity, error) {
	ret := _m.Called(entity)

	var r0 *GroupPermissionsEntity
	if rf, ok := ret.Get(0).(func(*GroupPermissionsEntity) *GroupPermissionsEntity); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GroupPermissionsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*GroupPermissionsEntity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type MockPermissionProvider struct {
	mock.Mock
}

// Find provides a mock function with given fields: criteria
func (_m *MockPermissionProvider) Find(criteria *PermissionCriteria) ([]*PermissionEntity, error) {
	ret := _m.Called(criteria)

	var r0 []*PermissionEntity
	if rf, ok := ret.Get(0).(func(*PermissionCriteria) []*PermissionEntity); ok {
		r0 = rf(criteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*PermissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*PermissionCriteria) error); ok {
		r1 = rf(criteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByID provides a mock function with given fields: id
func (_m *MockPermissionProvider) FindOneByID(id int64) (*PermissionEntity, error) {
	ret := _m.Called(id)

	var r0 *PermissionEntity
	if rf, ok := ret.Get(0).(func(int64) *PermissionEntity); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PermissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUserID provides a mock function with given fields: userID
func (_m *MockPermissionProvider) FindByUserID(userID int64) ([]*PermissionEntity, error) {
	ret := _m.Called(userID)

	var r0 []*PermissionEntity
	if rf, ok := ret.Get(0).(func(int64) []*PermissionEntity); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*PermissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByGroupID provides a mock function with given fields: groupID
func (_m *MockPermissionProvider) FindByGroupID(groupID int64) ([]*PermissionEntity, error) {
	ret := _m.Called(groupID)

	var r0 []*PermissionEntity
	if rf, ok := ret.Get(0).(func(int64) []*PermissionEntity); ok {
		r0 = rf(groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*PermissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: permissions
func (_m *MockPermissionProvider) Register(permissions charon.Permissions) (int64, int64, int64, error) {
	ret := _m.Called(permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(charon.Permissions) int64); ok {
		r0 = rf(permissions)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(charon.Permissions) int64); ok {
		r1 = rf(permissions)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 int64
	if rf, ok := ret.Get(2).(func(charon.Permissions) int64); ok {
		r2 = rf(permissions)
	} else {
		r2 = ret.Get(2).(int64)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(charon.Permissions) error); ok {
		r3 = rf(permissions)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// Insert provides a mock function with given fields: entity
func (_m *MockPermissionProvider) Insert(entity *PermissionEntity) (*PermissionEntity, error) {
	ret := _m.Called(entity)

	var r0 *PermissionEntity
	if rf, ok := ret.Get(0).(func(*PermissionEntity) *PermissionEntity); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PermissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*PermissionEntity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type MockPermissionRegistry struct {
	mock.Mock
}

// Exists provides a mock function with given fields: permission
func (_m *MockPermissionRegistry) Exists(permission charon.Permission) bool {
	ret := _m.Called(permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(charon.Permission) bool); ok {
		r0 = rf(permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Register provides a mock function with given fields: permissions
func (_m *MockPermissionRegistry) Register(permissions charon.Permissions) (int64, int64, int64, error) {
	ret := _m.Called(permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(charon.Permissions) int64); ok {
		r0 = rf(permissions)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(charon.Permissions) int64); ok {
		r1 = rf(permissions)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 int64
	if rf, ok := ret.Get(2).(func(charon.Permissions) int64); ok {
		r2 = rf(permissions)
	} else {
		r2 = ret.Get(2).(int64)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(charon.Permissions) error); ok {
		r3 = rf(permissions)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
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

// Exists provides a mock function with given fields: id
func (_m *MockUserProvider) Exists(id int64) (bool, error) {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int64) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: username, password, FirstName, LastName, confirmationToken, isSuperuser, IsStaff, isActive, isConfirmed
func (_m *MockUserProvider) Create(username string, password []byte, FirstName string, LastName string, confirmationToken []byte, isSuperuser bool, IsStaff bool, isActive bool, isConfirmed bool) (*UserEntity, error) {
	ret := _m.Called(username, password, FirstName, LastName, confirmationToken, isSuperuser, IsStaff, isActive, isConfirmed)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(string, []byte, string, string, []byte, bool, bool, bool, bool) *UserEntity); ok {
		r0 = rf(username, password, FirstName, LastName, confirmationToken, isSuperuser, IsStaff, isActive, isConfirmed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte, string, string, []byte, bool, bool, bool, bool) error); ok {
		r1 = rf(username, password, FirstName, LastName, confirmationToken, isSuperuser, IsStaff, isActive, isConfirmed)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: _a0
func (_m *MockUserProvider) Insert(_a0 *UserEntity) (*UserEntity, error) {
	ret := _m.Called(_a0)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(*UserEntity) *UserEntity); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*UserEntity) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSuperuser provides a mock function with given fields: username, password, FirstName, LastName
func (_m *MockUserProvider) CreateSuperuser(username string, password []byte, FirstName string, LastName string) (*UserEntity, error) {
	ret := _m.Called(username, password, FirstName, LastName)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(string, []byte, string, string) *UserEntity); ok {
		r0 = rf(username, password, FirstName, LastName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte, string, string) error); ok {
		r1 = rf(username, password, FirstName, LastName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Count provides a mock function with given fields:
func (_m *MockUserProvider) Count() (int64, error) {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateLastLoginAt provides a mock function with given fields: id
func (_m *MockUserProvider) UpdateLastLoginAt(id int64) (int64, error) {
	ret := _m.Called(id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64) int64); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChangePassword provides a mock function with given fields: id, password
func (_m *MockUserProvider) ChangePassword(id int64, password string) error {
	ret := _m.Called(id, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, string) error); ok {
		r0 = rf(id, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: criteria
func (_m *MockUserProvider) Find(criteria *UserCriteria) ([]*UserEntity, error) {
	ret := _m.Called(criteria)

	var r0 []*UserEntity
	if rf, ok := ret.Get(0).(func(*UserCriteria) []*UserEntity); ok {
		r0 = rf(criteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*UserCriteria) error); ok {
		r1 = rf(criteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByID provides a mock function with given fields: id
func (_m *MockUserProvider) FindOneByID(id int64) (*UserEntity, error) {
	ret := _m.Called(id)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(int64) *UserEntity); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByUsername provides a mock function with given fields: username
func (_m *MockUserProvider) FindOneByUsername(username string) (*UserEntity, error) {
	ret := _m.Called(username)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(string) *UserEntity); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOneByID provides a mock function with given fields: id
func (_m *MockUserProvider) DeleteOneByID(id int64) (int64, error) {
	ret := _m.Called(id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64) int64); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOneByID provides a mock function with given fields: _a0, _a1
func (_m *MockUserProvider) UpdateOneByID(_a0 int64, _a1 *UserPatch) (*UserEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *UserEntity
	if rf, ok := ret.Get(0).(func(int64, *UserPatch) *UserEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, *UserPatch) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegistrationConfirmation provides a mock function with given fields: id, confirmationToken
func (_m *MockUserProvider) RegistrationConfirmation(id int64, confirmationToken string) error {
	ret := _m.Called(id, confirmationToken)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, string) error); ok {
		r0 = rf(id, confirmationToken)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsGranted provides a mock function with given fields: id, permission
func (_m *MockUserProvider) IsGranted(id int64, permission charon.Permission) (bool, error) {
	ret := _m.Called(id, permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int64, charon.Permission) bool); ok {
		r0 = rf(id, permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, charon.Permission) error); ok {
		r1 = rf(id, permission)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPermissions provides a mock function with given fields: id, permissions
func (_m *MockUserProvider) SetPermissions(id int64, permissions ...charon.Permission) (int64, int64, error) {
	ret := _m.Called(id, permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, ...charon.Permission) int64); ok {
		r0 = rf(id, permissions...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(int64, ...charon.Permission) int64); ok {
		r1 = rf(id, permissions...)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int64, ...charon.Permission) error); ok {
		r2 = rf(id, permissions...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type MockUserGroupsProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: entity
func (_m *MockUserGroupsProvider) Insert(entity *UserGroupsEntity) (*UserGroupsEntity, error) {
	ret := _m.Called(entity)

	var r0 *UserGroupsEntity
	if rf, ok := ret.Get(0).(func(*UserGroupsEntity) *UserGroupsEntity); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserGroupsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*UserGroupsEntity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: userID, groupID
func (_m *MockUserGroupsProvider) Exists(userID int64, groupID int64) (bool, error) {
	ret := _m.Called(userID, groupID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int64, int64) bool); ok {
		r0 = rf(userID, groupID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(userID, groupID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: criteria
func (_m *MockUserGroupsProvider) Find(criteria *UserGroupsCriteria) ([]*UserGroupsEntity, error) {
	ret := _m.Called(criteria)

	var r0 []*UserGroupsEntity
	if rf, ok := ret.Get(0).(func(*UserGroupsCriteria) []*UserGroupsEntity); ok {
		r0 = rf(criteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*UserGroupsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*UserGroupsCriteria) error); ok {
		r1 = rf(criteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: userID, groupIDs
func (_m *MockUserGroupsProvider) Set(userID int64, groupIDs []int64) (int64, int64, error) {
	ret := _m.Called(userID, groupIDs)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, []int64) int64); ok {
		r0 = rf(userID, groupIDs)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(int64, []int64) int64); ok {
		r1 = rf(userID, groupIDs)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int64, []int64) error); ok {
		r2 = rf(userID, groupIDs)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteByUserID provides a mock function with given fields: id
func (_m *MockUserGroupsProvider) DeleteByUserID(id int64) (int64, error) {
	ret := _m.Called(id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64) int64); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type MockUserPermissionsProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: entity
func (_m *MockUserPermissionsProvider) Insert(entity *UserPermissionsEntity) (*UserPermissionsEntity, error) {
	ret := _m.Called(entity)

	var r0 *UserPermissionsEntity
	if rf, ok := ret.Get(0).(func(*UserPermissionsEntity) *UserPermissionsEntity); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*UserPermissionsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*UserPermissionsEntity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByUserID provides a mock function with given fields: id
func (_m *MockUserPermissionsProvider) DeleteByUserID(id int64) (int64, error) {
	ret := _m.Called(id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64) int64); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
