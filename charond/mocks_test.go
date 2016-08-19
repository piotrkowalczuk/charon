package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/stretchr/testify/mock"
)

import "github.com/piotrkowalczuk/ntypes"

type mockGroupProvider struct {
	mock.Mock
}

// insert provides a mock function with given fields: entity
func (_m *mockGroupProvider) insert(entity *groupEntity) (*groupEntity, error) {
	ret := _m.Called(entity)

	var r0 *groupEntity
	if rf, ok := ret.Get(0).(func(*groupEntity) *groupEntity); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*groupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*groupEntity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// findByUserID provides a mock function with given fields: _a0
func (_m *mockGroupProvider) findByUserID(_a0 int64) ([]*groupEntity, error) {
	ret := _m.Called(_a0)

	var r0 []*groupEntity
	if rf, ok := ret.Get(0).(func(int64) []*groupEntity); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*groupEntity)
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

// findOneByID provides a mock function with given fields: _a0
func (_m *mockGroupProvider) findOneByID(_a0 int64) (*groupEntity, error) {
	ret := _m.Called(_a0)

	var r0 *groupEntity
	if rf, ok := ret.Get(0).(func(int64) *groupEntity); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*groupEntity)
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

// find provides a mock function with given fields: c
func (_m *mockGroupProvider) find(c *groupCriteria) ([]*groupEntity, error) {
	ret := _m.Called(c)

	var r0 []*groupEntity
	if rf, ok := ret.Get(0).(func(*groupCriteria) []*groupEntity); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*groupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*groupCriteria) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// create provides a mock function with given fields: createdBy, name, description
func (_m *mockGroupProvider) create(createdBy int64, name string, description *ntypes.String) (*groupEntity, error) {
	ret := _m.Called(createdBy, name, description)

	var r0 *groupEntity
	if rf, ok := ret.Get(0).(func(int64, string, *ntypes.String) *groupEntity); ok {
		r0 = rf(createdBy, name, description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*groupEntity)
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

// updateOneByID provides a mock function with given fields: id, updatedBy, name, description
func (_m *mockGroupProvider) updateOneByID(id int64, updatedBy int64, name *ntypes.String, description *ntypes.String) (*groupEntity, error) {
	ret := _m.Called(id, updatedBy, name, description)

	var r0 *groupEntity
	if rf, ok := ret.Get(0).(func(int64, int64, *ntypes.String, *ntypes.String) *groupEntity); ok {
		r0 = rf(id, updatedBy, name, description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*groupEntity)
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

// deleteOneByID provides a mock function with given fields: id
func (_m *mockGroupProvider) deleteOneByID(id int64) (int64, error) {
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

// isGranted provides a mock function with given fields: id, permission
func (_m *mockGroupProvider) isGranted(id int64, permission charon.Permission) (bool, error) {
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

// setPermissions provides a mock function with given fields: id, permissions
func (_m *mockGroupProvider) setPermissions(id int64, permissions ...charon.Permission) (int64, int64, error) {
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

type mockGroupPermissionsProvider struct {
	mock.Mock
}

// insert provides a mock function with given fields: entity
func (_m *mockGroupPermissionsProvider) insert(entity *groupPermissionsEntity) (*groupPermissionsEntity, error) {
	ret := _m.Called(entity)

	var r0 *groupPermissionsEntity
	if rf, ok := ret.Get(0).(func(*groupPermissionsEntity) *groupPermissionsEntity); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*groupPermissionsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*groupPermissionsEntity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockPermissionProvider struct {
	mock.Mock
}

// find provides a mock function with given fields: criteria
func (_m *mockPermissionProvider) find(criteria *permissionCriteria) ([]*permissionEntity, error) {
	ret := _m.Called(criteria)

	var r0 []*permissionEntity
	if rf, ok := ret.Get(0).(func(*permissionCriteria) []*permissionEntity); ok {
		r0 = rf(criteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*permissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*permissionCriteria) error); ok {
		r1 = rf(criteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// findOneByID provides a mock function with given fields: id
func (_m *mockPermissionProvider) findOneByID(id int64) (*permissionEntity, error) {
	ret := _m.Called(id)

	var r0 *permissionEntity
	if rf, ok := ret.Get(0).(func(int64) *permissionEntity); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*permissionEntity)
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

// findByUserID provides a mock function with given fields: userID
func (_m *mockPermissionProvider) findByUserID(userID int64) ([]*permissionEntity, error) {
	ret := _m.Called(userID)

	var r0 []*permissionEntity
	if rf, ok := ret.Get(0).(func(int64) []*permissionEntity); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*permissionEntity)
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

// findByGroupID provides a mock function with given fields: groupID
func (_m *mockPermissionProvider) findByGroupID(groupID int64) ([]*permissionEntity, error) {
	ret := _m.Called(groupID)

	var r0 []*permissionEntity
	if rf, ok := ret.Get(0).(func(int64) []*permissionEntity); ok {
		r0 = rf(groupID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*permissionEntity)
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

// register provides a mock function with given fields: permissions
func (_m *mockPermissionProvider) register(permissions charon.Permissions) (int64, int64, int64, error) {
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

// insert provides a mock function with given fields: entity
func (_m *mockPermissionProvider) insert(entity *permissionEntity) (*permissionEntity, error) {
	ret := _m.Called(entity)

	var r0 *permissionEntity
	if rf, ok := ret.Get(0).(func(*permissionEntity) *permissionEntity); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*permissionEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*permissionEntity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockPermissionRegistry struct {
	mock.Mock
}

// exists provides a mock function with given fields: permission
func (_m *mockPermissionRegistry) exists(permission charon.Permission) bool {
	ret := _m.Called(permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(charon.Permission) bool); ok {
		r0 = rf(permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// register provides a mock function with given fields: permissions
func (_m *mockPermissionRegistry) register(permissions charon.Permissions) (int64, int64, int64, error) {
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

type mockUserProvider struct {
	mock.Mock
}

// Exists provides a mock function with given fields: id
func (_m *mockUserProvider) Exists(id int64) (bool, error) {
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

// Create provides a mock function with given fields: username, password, firstName, lastName, confirmationToken, isSuperuser, isStaff, isActive, isConfirmed
func (_m *mockUserProvider) Create(username string, password []byte, firstName string, lastName string, confirmationToken []byte, isSuperuser bool, isStaff bool, isActive bool, isConfirmed bool) (*userEntity, error) {
	ret := _m.Called(username, password, firstName, lastName, confirmationToken, isSuperuser, isStaff, isActive, isConfirmed)

	var r0 *userEntity
	if rf, ok := ret.Get(0).(func(string, []byte, string, string, []byte, bool, bool, bool, bool) *userEntity); ok {
		r0 = rf(username, password, firstName, lastName, confirmationToken, isSuperuser, isStaff, isActive, isConfirmed)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte, string, string, []byte, bool, bool, bool, bool) error); ok {
		r1 = rf(username, password, firstName, lastName, confirmationToken, isSuperuser, isStaff, isActive, isConfirmed)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// insert provides a mock function with given fields: _a0
func (_m *mockUserProvider) insert(_a0 *userEntity) (*userEntity, error) {
	ret := _m.Called(_a0)

	var r0 *userEntity
	if rf, ok := ret.Get(0).(func(*userEntity) *userEntity); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*userEntity) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSuperuser provides a mock function with given fields: username, password, firstName, lastName
func (_m *mockUserProvider) CreateSuperuser(username string, password []byte, firstName string, lastName string) (*userEntity, error) {
	ret := _m.Called(username, password, firstName, lastName)

	var r0 *userEntity
	if rf, ok := ret.Get(0).(func(string, []byte, string, string) *userEntity); ok {
		r0 = rf(username, password, firstName, lastName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte, string, string) error); ok {
		r1 = rf(username, password, firstName, lastName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Count provides a mock function with given fields:
func (_m *mockUserProvider) Count() (int64, error) {
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

// updateLastLoginAt provides a mock function with given fields: id
func (_m *mockUserProvider) updateLastLoginAt(id int64) (int64, error) {
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
func (_m *mockUserProvider) ChangePassword(id int64, password string) error {
	ret := _m.Called(id, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, string) error); ok {
		r0 = rf(id, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// find provides a mock function with given fields: criteria
func (_m *mockUserProvider) find(criteria *userCriteria) ([]*userEntity, error) {
	ret := _m.Called(criteria)

	var r0 []*userEntity
	if rf, ok := ret.Get(0).(func(*userCriteria) []*userEntity); ok {
		r0 = rf(criteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*userEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*userCriteria) error); ok {
		r1 = rf(criteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// findOneByID provides a mock function with given fields: id
func (_m *mockUserProvider) findOneByID(id int64) (*userEntity, error) {
	ret := _m.Called(id)

	var r0 *userEntity
	if rf, ok := ret.Get(0).(func(int64) *userEntity); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userEntity)
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

// findOneByUsername provides a mock function with given fields: username
func (_m *mockUserProvider) findOneByUsername(username string) (*userEntity, error) {
	ret := _m.Called(username)

	var r0 *userEntity
	if rf, ok := ret.Get(0).(func(string) *userEntity); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userEntity)
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

// deleteOneByID provides a mock function with given fields: id
func (_m *mockUserProvider) deleteOneByID(id int64) (int64, error) {
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

// updateOneByID provides a mock function with given fields: _a0, _a1
func (_m *mockUserProvider) updateOneByID(_a0 int64, _a1 *userPatch) (*userEntity, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userEntity
	if rf, ok := ret.Get(0).(func(int64, *userPatch) *userEntity); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, *userPatch) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegistrationConfirmation provides a mock function with given fields: id, confirmationToken
func (_m *mockUserProvider) RegistrationConfirmation(id int64, confirmationToken string) error {
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
func (_m *mockUserProvider) IsGranted(id int64, permission charon.Permission) (bool, error) {
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
func (_m *mockUserProvider) SetPermissions(id int64, permissions ...charon.Permission) (int64, int64, error) {
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

type mockUserGroupsProvider struct {
	mock.Mock
}

// insert provides a mock function with given fields: entity
func (_m *mockUserGroupsProvider) insert(entity *userGroupsEntity) (*userGroupsEntity, error) {
	ret := _m.Called(entity)

	var r0 *userGroupsEntity
	if rf, ok := ret.Get(0).(func(*userGroupsEntity) *userGroupsEntity); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGroupsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*userGroupsEntity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exists provides a mock function with given fields: userID, groupID
func (_m *mockUserGroupsProvider) Exists(userID int64, groupID int64) (bool, error) {
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

// find provides a mock function with given fields: criteria
func (_m *mockUserGroupsProvider) find(criteria *userGroupsCriteria) ([]*userGroupsEntity, error) {
	ret := _m.Called(criteria)

	var r0 []*userGroupsEntity
	if rf, ok := ret.Get(0).(func(*userGroupsCriteria) []*userGroupsEntity); ok {
		r0 = rf(criteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*userGroupsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*userGroupsCriteria) error); ok {
		r1 = rf(criteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: userID, groupIDs
func (_m *mockUserGroupsProvider) Set(userID int64, groupIDs []int64) (int64, int64, error) {
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

type mockUserPermissionsProvider struct {
	mock.Mock
}

// insert provides a mock function with given fields: entity
func (_m *mockUserPermissionsProvider) insert(entity *userPermissionsEntity) (*userPermissionsEntity, error) {
	ret := _m.Called(entity)

	var r0 *userPermissionsEntity
	if rf, ok := ret.Get(0).(func(*userPermissionsEntity) *userPermissionsEntity); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userPermissionsEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*userPermissionsEntity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
