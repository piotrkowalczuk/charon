package charon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

import (
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/nilt"
)
import "golang.org/x/net/context"
import "google.golang.org/grpc"

type MockCharon struct {
	mock.Mock
}

// IsGranted provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockCharon) IsGranted(_a0 context.Context, _a1 int64, _a2 Permission) (bool, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64, Permission) bool); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, Permission) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthenticated provides a mock function with given fields: _a0, _a1
func (_m *MockCharon) IsAuthenticated(_a0 context.Context, _a1 mnemosyne.AccessToken) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, mnemosyne.AccessToken) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, mnemosyne.AccessToken) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subject provides a mock function with given fields: _a0, _a1
func (_m *MockCharon) Subject(_a0 context.Context, _a1 mnemosyne.AccessToken) (*Subject, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *Subject
	if rf, ok := ret.Get(0).(func(context.Context, mnemosyne.AccessToken) *Subject); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Subject)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, mnemosyne.AccessToken) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FromContext provides a mock function with given fields: _a0
func (_m *MockCharon) FromContext(_a0 context.Context) (*Subject, error) {
	ret := _m.Called(_a0)

	var r0 *Subject
	if rf, ok := ret.Get(0).(func(context.Context) *Subject); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Subject)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockCharon) Login(_a0 context.Context, _a1 string, _a2 string) (*mnemosyne.AccessToken, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *mnemosyne.AccessToken
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *mnemosyne.AccessToken); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mnemosyne.AccessToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: _a0, _a1
func (_m *MockCharon) Logout(_a0 context.Context, _a1 mnemosyne.AccessToken) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, mnemosyne.AccessToken) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type MockRPCClient struct {
	mock.Mock
}

// Login provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *LoginResponse
	if rf, ok := ret.Get(0).(func(context.Context, *LoginRequest, ...grpc.CallOption) *LoginResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*LoginResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *LoginRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *LogoutResponse
	if rf, ok := ret.Get(0).(func(context.Context, *LogoutRequest, ...grpc.CallOption) *LogoutResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*LogoutResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *LogoutRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthenticated provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) IsAuthenticated(ctx context.Context, in *IsAuthenticatedRequest, opts ...grpc.CallOption) (*IsAuthenticatedResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *IsAuthenticatedResponse
	if rf, ok := ret.Get(0).(func(context.Context, *IsAuthenticatedRequest, ...grpc.CallOption) *IsAuthenticatedResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*IsAuthenticatedResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *IsAuthenticatedRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subject provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) Subject(ctx context.Context, in *SubjectRequest, opts ...grpc.CallOption) (*SubjectResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *SubjectResponse
	if rf, ok := ret.Get(0).(func(context.Context, *SubjectRequest, ...grpc.CallOption) *SubjectResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*SubjectResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *SubjectRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsGranted provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) IsGranted(ctx context.Context, in *IsGrantedRequest, opts ...grpc.CallOption) (*IsGrantedResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *IsGrantedResponse
	if rf, ok := ret.Get(0).(func(context.Context, *IsGrantedRequest, ...grpc.CallOption) *IsGrantedResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*IsGrantedResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *IsGrantedRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BelongsTo provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) BelongsTo(ctx context.Context, in *BelongsToRequest, opts ...grpc.CallOption) (*BelongsToResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *BelongsToResponse
	if rf, ok := ret.Get(0).(func(context.Context, *BelongsToRequest, ...grpc.CallOption) *BelongsToResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*BelongsToResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *BelongsToRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *CreateUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *CreateUserRequest, ...grpc.CallOption) *CreateUserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*CreateUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *CreateUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModifyUser provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) ModifyUser(ctx context.Context, in *ModifyUserRequest, opts ...grpc.CallOption) (*ModifyUserResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *ModifyUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ModifyUserRequest, ...grpc.CallOption) *ModifyUserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ModifyUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ModifyUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *GetUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *GetUserRequest, ...grpc.CallOption) *GetUserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GetUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *GetUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUsers provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *ListUsersResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListUsersRequest, ...grpc.CallOption) *ListUsersResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListUsersResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListUsersRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *DeleteUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *DeleteUserRequest, ...grpc.CallOption) *DeleteUserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DeleteUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *DeleteUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUserPermissions provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) ListUserPermissions(ctx context.Context, in *ListUserPermissionsRequest, opts ...grpc.CallOption) (*ListUserPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *ListUserPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListUserPermissionsRequest, ...grpc.CallOption) *ListUserPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListUserPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListUserPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetUserPermissions provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) SetUserPermissions(ctx context.Context, in *SetUserPermissionsRequest, opts ...grpc.CallOption) (*SetUserPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *SetUserPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *SetUserPermissionsRequest, ...grpc.CallOption) *SetUserPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*SetUserPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *SetUserPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUserGroups provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) ListUserGroups(ctx context.Context, in *ListUserGroupsRequest, opts ...grpc.CallOption) (*ListUserGroupsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *ListUserGroupsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListUserGroupsRequest, ...grpc.CallOption) *ListUserGroupsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListUserGroupsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListUserGroupsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetUserGroups provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) SetUserGroups(ctx context.Context, in *SetUserGroupsRequest, opts ...grpc.CallOption) (*SetUserGroupsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *SetUserGroupsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *SetUserGroupsRequest, ...grpc.CallOption) *SetUserGroupsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*SetUserGroupsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *SetUserGroupsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterPermissions provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) RegisterPermissions(ctx context.Context, in *RegisterPermissionsRequest, opts ...grpc.CallOption) (*RegisterPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *RegisterPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *RegisterPermissionsRequest, ...grpc.CallOption) *RegisterPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*RegisterPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *RegisterPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPermissions provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) ListPermissions(ctx context.Context, in *ListPermissionsRequest, opts ...grpc.CallOption) (*ListPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *ListPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListPermissionsRequest, ...grpc.CallOption) *ListPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPermission provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) GetPermission(ctx context.Context, in *GetPermissionRequest, opts ...grpc.CallOption) (*GetPermissionResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *GetPermissionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *GetPermissionRequest, ...grpc.CallOption) *GetPermissionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GetPermissionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *GetPermissionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateGroup provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) CreateGroup(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*CreateGroupResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *CreateGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *CreateGroupRequest, ...grpc.CallOption) *CreateGroupResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*CreateGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *CreateGroupRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModifyGroup provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) ModifyGroup(ctx context.Context, in *ModifyGroupRequest, opts ...grpc.CallOption) (*ModifyGroupResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *ModifyGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ModifyGroupRequest, ...grpc.CallOption) *ModifyGroupResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ModifyGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ModifyGroupRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGroup provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) GetGroup(ctx context.Context, in *GetGroupRequest, opts ...grpc.CallOption) (*GetGroupResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *GetGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *GetGroupRequest, ...grpc.CallOption) *GetGroupResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GetGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *GetGroupRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGroups provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) ListGroups(ctx context.Context, in *ListGroupsRequest, opts ...grpc.CallOption) (*ListGroupsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *ListGroupsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListGroupsRequest, ...grpc.CallOption) *ListGroupsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListGroupsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListGroupsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGroup provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) DeleteGroup(ctx context.Context, in *DeleteGroupRequest, opts ...grpc.CallOption) (*DeleteGroupResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *DeleteGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *DeleteGroupRequest, ...grpc.CallOption) *DeleteGroupResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DeleteGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *DeleteGroupRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGroupPermissions provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) ListGroupPermissions(ctx context.Context, in *ListGroupPermissionsRequest, opts ...grpc.CallOption) (*ListGroupPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *ListGroupPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListGroupPermissionsRequest, ...grpc.CallOption) *ListGroupPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListGroupPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListGroupPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetGroupPermissions provides a mock function with given fields: ctx, in, opts
func (_m *MockRPCClient) SetGroupPermissions(ctx context.Context, in *SetGroupPermissionsRequest, opts ...grpc.CallOption) (*SetGroupPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *SetGroupPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *SetGroupPermissionsRequest, ...grpc.CallOption) *SetGroupPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*SetGroupPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *SetGroupPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type MockRPCServer struct {
	mock.Mock
}

// Login provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) Login(_a0 context.Context, _a1 *LoginRequest) (*LoginResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *LoginResponse
	if rf, ok := ret.Get(0).(func(context.Context, *LoginRequest) *LoginResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*LoginResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *LoginRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) Logout(_a0 context.Context, _a1 *LogoutRequest) (*LogoutResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *LogoutResponse
	if rf, ok := ret.Get(0).(func(context.Context, *LogoutRequest) *LogoutResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*LogoutResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *LogoutRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthenticated provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) IsAuthenticated(_a0 context.Context, _a1 *IsAuthenticatedRequest) (*IsAuthenticatedResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *IsAuthenticatedResponse
	if rf, ok := ret.Get(0).(func(context.Context, *IsAuthenticatedRequest) *IsAuthenticatedResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*IsAuthenticatedResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *IsAuthenticatedRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subject provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) Subject(_a0 context.Context, _a1 *SubjectRequest) (*SubjectResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *SubjectResponse
	if rf, ok := ret.Get(0).(func(context.Context, *SubjectRequest) *SubjectResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*SubjectResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *SubjectRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsGranted provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) IsGranted(_a0 context.Context, _a1 *IsGrantedRequest) (*IsGrantedResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *IsGrantedResponse
	if rf, ok := ret.Get(0).(func(context.Context, *IsGrantedRequest) *IsGrantedResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*IsGrantedResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *IsGrantedRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BelongsTo provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) BelongsTo(_a0 context.Context, _a1 *BelongsToRequest) (*BelongsToResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *BelongsToResponse
	if rf, ok := ret.Get(0).(func(context.Context, *BelongsToRequest) *BelongsToResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*BelongsToResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *BelongsToRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) CreateUser(_a0 context.Context, _a1 *CreateUserRequest) (*CreateUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *CreateUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *CreateUserRequest) *CreateUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*CreateUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *CreateUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModifyUser provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) ModifyUser(_a0 context.Context, _a1 *ModifyUserRequest) (*ModifyUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ModifyUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ModifyUserRequest) *ModifyUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ModifyUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ModifyUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) GetUser(_a0 context.Context, _a1 *GetUserRequest) (*GetUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *GetUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *GetUserRequest) *GetUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GetUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *GetUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUsers provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) ListUsers(_a0 context.Context, _a1 *ListUsersRequest) (*ListUsersResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ListUsersResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListUsersRequest) *ListUsersResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListUsersResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListUsersRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) DeleteUser(_a0 context.Context, _a1 *DeleteUserRequest) (*DeleteUserResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *DeleteUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *DeleteUserRequest) *DeleteUserResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DeleteUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *DeleteUserRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUserPermissions provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) ListUserPermissions(_a0 context.Context, _a1 *ListUserPermissionsRequest) (*ListUserPermissionsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ListUserPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListUserPermissionsRequest) *ListUserPermissionsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListUserPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListUserPermissionsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetUserPermissions provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) SetUserPermissions(_a0 context.Context, _a1 *SetUserPermissionsRequest) (*SetUserPermissionsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *SetUserPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *SetUserPermissionsRequest) *SetUserPermissionsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*SetUserPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *SetUserPermissionsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUserGroups provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) ListUserGroups(_a0 context.Context, _a1 *ListUserGroupsRequest) (*ListUserGroupsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ListUserGroupsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListUserGroupsRequest) *ListUserGroupsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListUserGroupsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListUserGroupsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetUserGroups provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) SetUserGroups(_a0 context.Context, _a1 *SetUserGroupsRequest) (*SetUserGroupsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *SetUserGroupsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *SetUserGroupsRequest) *SetUserGroupsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*SetUserGroupsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *SetUserGroupsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterPermissions provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) RegisterPermissions(_a0 context.Context, _a1 *RegisterPermissionsRequest) (*RegisterPermissionsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *RegisterPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *RegisterPermissionsRequest) *RegisterPermissionsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*RegisterPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *RegisterPermissionsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPermissions provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) ListPermissions(_a0 context.Context, _a1 *ListPermissionsRequest) (*ListPermissionsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ListPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListPermissionsRequest) *ListPermissionsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListPermissionsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPermission provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) GetPermission(_a0 context.Context, _a1 *GetPermissionRequest) (*GetPermissionResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *GetPermissionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *GetPermissionRequest) *GetPermissionResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GetPermissionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *GetPermissionRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateGroup provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) CreateGroup(_a0 context.Context, _a1 *CreateGroupRequest) (*CreateGroupResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *CreateGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *CreateGroupRequest) *CreateGroupResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*CreateGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *CreateGroupRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModifyGroup provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) ModifyGroup(_a0 context.Context, _a1 *ModifyGroupRequest) (*ModifyGroupResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ModifyGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ModifyGroupRequest) *ModifyGroupResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ModifyGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ModifyGroupRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGroup provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) GetGroup(_a0 context.Context, _a1 *GetGroupRequest) (*GetGroupResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *GetGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *GetGroupRequest) *GetGroupResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GetGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *GetGroupRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGroups provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) ListGroups(_a0 context.Context, _a1 *ListGroupsRequest) (*ListGroupsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ListGroupsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListGroupsRequest) *ListGroupsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListGroupsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListGroupsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGroup provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) DeleteGroup(_a0 context.Context, _a1 *DeleteGroupRequest) (*DeleteGroupResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *DeleteGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *DeleteGroupRequest) *DeleteGroupResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DeleteGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *DeleteGroupRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGroupPermissions provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) ListGroupPermissions(_a0 context.Context, _a1 *ListGroupPermissionsRequest) (*ListGroupPermissionsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ListGroupPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *ListGroupPermissionsRequest) *ListGroupPermissionsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListGroupPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListGroupPermissionsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetGroupPermissions provides a mock function with given fields: _a0, _a1
func (_m *MockRPCServer) SetGroupPermissions(_a0 context.Context, _a1 *SetGroupPermissionsRequest) (*SetGroupPermissionsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *SetGroupPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *SetGroupPermissionsRequest) *SetGroupPermissionsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*SetGroupPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *SetGroupPermissionsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockGroupProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: entity
func (_m *mockGroupProvider) Insert(entity *groupEntity) (*groupEntity, error) {
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

// FindByUserID provides a mock function with given fields: _a0
func (_m *mockGroupProvider) FindByUserID(_a0 int64) ([]*groupEntity, error) {
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

// FindOneByID provides a mock function with given fields: _a0
func (_m *mockGroupProvider) FindOneByID(_a0 int64) (*groupEntity, error) {
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

// Find provides a mock function with given fields: c
func (_m *mockGroupProvider) Find(c *groupCriteria) ([]*groupEntity, error) {
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

// Create provides a mock function with given fields: createdBy, name, description
func (_m *mockGroupProvider) Create(createdBy int64, name string, description *nilt.String) (*groupEntity, error) {
	ret := _m.Called(createdBy, name, description)

	var r0 *groupEntity
	if rf, ok := ret.Get(0).(func(int64, string, *nilt.String) *groupEntity); ok {
		r0 = rf(createdBy, name, description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*groupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, string, *nilt.String) error); ok {
		r1 = rf(createdBy, name, description)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOneByID provides a mock function with given fields: id, updatedBy, name, description
func (_m *mockGroupProvider) UpdateOneByID(id int64, updatedBy int64, name *nilt.String, description *nilt.String) (*groupEntity, error) {
	ret := _m.Called(id, updatedBy, name, description)

	var r0 *groupEntity
	if rf, ok := ret.Get(0).(func(int64, int64, *nilt.String, *nilt.String) *groupEntity); ok {
		r0 = rf(id, updatedBy, name, description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*groupEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64, *nilt.String, *nilt.String) error); ok {
		r1 = rf(id, updatedBy, name, description)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByID provides a mock function with given fields: id
func (_m *mockGroupProvider) DeleteByID(id int64) (int64, error) {
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
func (_m *mockGroupProvider) IsGranted(id int64, permission Permission) (bool, error) {
	ret := _m.Called(id, permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int64, Permission) bool); ok {
		r0 = rf(id, permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, Permission) error); ok {
		r1 = rf(id, permission)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPermissions provides a mock function with given fields: id, permissions
func (_m *mockGroupProvider) SetPermissions(id int64, permissions ...Permission) (int64, int64, error) {
	ret := _m.Called(id, permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, ...Permission) int64); ok {
		r0 = rf(id, permissions...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(int64, ...Permission) int64); ok {
		r1 = rf(id, permissions...)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int64, ...Permission) error); ok {
		r2 = rf(id, permissions...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockGroupPermissionsProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: entity
func (_m *mockGroupPermissionsProvider) Insert(entity *groupPermissionsEntity) (*groupPermissionsEntity, error) {
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

type MockPasswordHasher struct {
	mock.Mock
}

// Hash provides a mock function with given fields: _a0
func (_m *MockPasswordHasher) Hash(_a0 []byte) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Compare provides a mock function with given fields: _a0, _a1
func (_m *MockPasswordHasher) Compare(_a0 []byte, _a1 []byte) bool {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func([]byte, []byte) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockPermissionProvider struct {
	mock.Mock
}

// Find provides a mock function with given fields: criteria
func (_m *mockPermissionProvider) Find(criteria *permissionCriteria) ([]*permissionEntity, error) {
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

// FindOneByID provides a mock function with given fields: id
func (_m *mockPermissionProvider) FindOneByID(id int64) (*permissionEntity, error) {
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

// FindByUserID provides a mock function with given fields: userID
func (_m *mockPermissionProvider) FindByUserID(userID int64) ([]*permissionEntity, error) {
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

// FindByGroupID provides a mock function with given fields: groupID
func (_m *mockPermissionProvider) FindByGroupID(groupID int64) ([]*permissionEntity, error) {
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

// Register provides a mock function with given fields: permissions
func (_m *mockPermissionProvider) Register(permissions Permissions) (int64, int64, int64, error) {
	ret := _m.Called(permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(Permissions) int64); ok {
		r0 = rf(permissions)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(Permissions) int64); ok {
		r1 = rf(permissions)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 int64
	if rf, ok := ret.Get(2).(func(Permissions) int64); ok {
		r2 = rf(permissions)
	} else {
		r2 = ret.Get(2).(int64)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(Permissions) error); ok {
		r3 = rf(permissions)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// Insert provides a mock function with given fields: entity
func (_m *mockPermissionProvider) Insert(entity *permissionEntity) (*permissionEntity, error) {
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

type MockPermissionRegistry struct {
	mock.Mock
}

// Exists provides a mock function with given fields: permission
func (_m *MockPermissionRegistry) Exists(permission Permission) bool {
	ret := _m.Called(permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(Permission) bool); ok {
		r0 = rf(permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Register provides a mock function with given fields: permissions
func (_m *MockPermissionRegistry) Register(permissions Permissions) (int64, int64, int64, error) {
	ret := _m.Called(permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(Permissions) int64); ok {
		r0 = rf(permissions)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(Permissions) int64); ok {
		r1 = rf(permissions)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 int64
	if rf, ok := ret.Get(2).(func(Permissions) int64); ok {
		r2 = rf(permissions)
	} else {
		r2 = ret.Get(2).(int64)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(Permissions) error); ok {
		r3 = rf(permissions)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

type MockSecurityContext struct {
	mock.Mock
}

// Subject provides a mock function with given fields:
func (_m *MockSecurityContext) Subject() (Subject, bool) {
	ret := _m.Called()

	var r0 Subject
	if rf, ok := ret.Get(0).(func() Subject); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(Subject)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// AccessToken provides a mock function with given fields:
func (_m *MockSecurityContext) AccessToken() (mnemosyne.AccessToken, bool) {
	ret := _m.Called()

	var r0 mnemosyne.AccessToken
	if rf, ok := ret.Get(0).(func() mnemosyne.AccessToken); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(mnemosyne.AccessToken)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
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

// Insert provides a mock function with given fields: _a0
func (_m *mockUserProvider) Insert(_a0 *userEntity) (*userEntity, error) {
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

// UpdateLastLoginAt provides a mock function with given fields: id
func (_m *mockUserProvider) UpdateLastLoginAt(id int64) (int64, error) {
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

// Find provides a mock function with given fields: criteria
func (_m *mockUserProvider) Find(criteria *userCriteria) ([]*userEntity, error) {
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

// FindOneByID provides a mock function with given fields: id
func (_m *mockUserProvider) FindOneByID(id int64) (*userEntity, error) {
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

// FindOneByUsername provides a mock function with given fields: username
func (_m *mockUserProvider) FindOneByUsername(username string) (*userEntity, error) {
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

// DeleteByID provides a mock function with given fields: id
func (_m *mockUserProvider) DeleteByID(id int64) (int64, error) {
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

// UpdateByID provides a mock function with given fields: id, confirmationToken, createdAt, createdBy, firstName, isActive, isConfirmed, isStaff, isSuperuser, lastLoginAt, lastName, password, updatedAt, updatedBy, username
func (_m *mockUserProvider) UpdateByID(id int64, confirmationToken []byte, createdAt *time.Time, createdBy *nilt.Int64, firstName *nilt.String, isActive *nilt.Bool, isConfirmed *nilt.Bool, isStaff *nilt.Bool, isSuperuser *nilt.Bool, lastLoginAt *time.Time, lastName *nilt.String, password []byte, updatedAt *time.Time, updatedBy *nilt.Int64, username *nilt.String) (*userEntity, error) {
	ret := _m.Called(id, confirmationToken, createdAt, createdBy, firstName, isActive, isConfirmed, isStaff, isSuperuser, lastLoginAt, lastName, password, updatedAt, updatedBy, username)

	var r0 *userEntity
	if rf, ok := ret.Get(0).(func(int64, []byte, *time.Time, *nilt.Int64, *nilt.String, *nilt.Bool, *nilt.Bool, *nilt.Bool, *nilt.Bool, *time.Time, *nilt.String, []byte, *time.Time, *nilt.Int64, *nilt.String) *userEntity); ok {
		r0 = rf(id, confirmationToken, createdAt, createdBy, firstName, isActive, isConfirmed, isStaff, isSuperuser, lastLoginAt, lastName, password, updatedAt, updatedBy, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, []byte, *time.Time, *nilt.Int64, *nilt.String, *nilt.Bool, *nilt.Bool, *nilt.Bool, *nilt.Bool, *time.Time, *nilt.String, []byte, *time.Time, *nilt.Int64, *nilt.String) error); ok {
		r1 = rf(id, confirmationToken, createdAt, createdBy, firstName, isActive, isConfirmed, isStaff, isSuperuser, lastLoginAt, lastName, password, updatedAt, updatedBy, username)
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
func (_m *mockUserProvider) IsGranted(id int64, permission Permission) (bool, error) {
	ret := _m.Called(id, permission)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int64, Permission) bool); ok {
		r0 = rf(id, permission)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, Permission) error); ok {
		r1 = rf(id, permission)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPermissions provides a mock function with given fields: id, permissions
func (_m *mockUserProvider) SetPermissions(id int64, permissions ...Permission) (int64, int64, error) {
	ret := _m.Called(id, permissions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, ...Permission) int64); ok {
		r0 = rf(id, permissions...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(int64, ...Permission) int64); ok {
		r1 = rf(id, permissions...)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int64, ...Permission) error); ok {
		r2 = rf(id, permissions...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockUserGroupsProvider struct {
	mock.Mock
}

// Insert provides a mock function with given fields: entity
func (_m *mockUserGroupsProvider) Insert(entity *userGroupsEntity) (*userGroupsEntity, error) {
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

// Find provides a mock function with given fields: criteria
func (_m *mockUserGroupsProvider) Find(criteria *userGroupsCriteria) ([]*userGroupsEntity, error) {
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

// Insert provides a mock function with given fields: entity
func (_m *mockUserPermissionsProvider) Insert(entity *userPermissionsEntity) (*userPermissionsEntity, error) {
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
