package charontest

import "github.com/piotrkowalczuk/charon/charonrpc"
import "github.com/stretchr/testify/mock"

import google_protobuf "github.com/golang/protobuf/ptypes/empty"
import google_protobuf1 "github.com/golang/protobuf/ptypes/wrappers"
import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"

type AuthClient struct {
	mock.Mock
}

// Login provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) Login(ctx context.Context, in *charonrpc.LoginRequest, opts ...grpc.CallOption) (*google_protobuf1.StringValue, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *google_protobuf1.StringValue
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.LoginRequest, ...grpc.CallOption) *google_protobuf1.StringValue); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*google_protobuf1.StringValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.LoginRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) Logout(ctx context.Context, in *charonrpc.LogoutRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *google_protobuf.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.LogoutRequest, ...grpc.CallOption) *google_protobuf.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*google_protobuf.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.LogoutRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthenticated provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) IsAuthenticated(ctx context.Context, in *charonrpc.IsAuthenticatedRequest, opts ...grpc.CallOption) (*google_protobuf1.BoolValue, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *google_protobuf1.BoolValue
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.IsAuthenticatedRequest, ...grpc.CallOption) *google_protobuf1.BoolValue); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*google_protobuf1.BoolValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.IsAuthenticatedRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Actor provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) Actor(ctx context.Context, in *google_protobuf1.StringValue, opts ...grpc.CallOption) (*charonrpc.ActorResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.ActorResponse
	if rf, ok := ret.Get(0).(func(context.Context, *google_protobuf1.StringValue, ...grpc.CallOption) *charonrpc.ActorResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.ActorResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *google_protobuf1.StringValue, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsGranted provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) IsGranted(ctx context.Context, in *charonrpc.IsGrantedRequest, opts ...grpc.CallOption) (*google_protobuf1.BoolValue, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *google_protobuf1.BoolValue
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.IsGrantedRequest, ...grpc.CallOption) *google_protobuf1.BoolValue); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*google_protobuf1.BoolValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.IsGrantedRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BelongsTo provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) BelongsTo(ctx context.Context, in *charonrpc.BelongsToRequest, opts ...grpc.CallOption) (*google_protobuf1.BoolValue, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *google_protobuf1.BoolValue
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.BelongsToRequest, ...grpc.CallOption) *google_protobuf1.BoolValue); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*google_protobuf1.BoolValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.BelongsToRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type GroupManagerClient struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in, opts
func (_m *GroupManagerClient) Create(ctx context.Context, in *charonrpc.CreateGroupRequest, opts ...grpc.CallOption) (*charonrpc.CreateGroupResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.CreateGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.CreateGroupRequest, ...grpc.CallOption) *charonrpc.CreateGroupResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.CreateGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.CreateGroupRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Modify provides a mock function with given fields: ctx, in, opts
func (_m *GroupManagerClient) Modify(ctx context.Context, in *charonrpc.ModifyGroupRequest, opts ...grpc.CallOption) (*charonrpc.ModifyGroupResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.ModifyGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.ModifyGroupRequest, ...grpc.CallOption) *charonrpc.ModifyGroupResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.ModifyGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.ModifyGroupRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *GroupManagerClient) Get(ctx context.Context, in *charonrpc.GetGroupRequest, opts ...grpc.CallOption) (*charonrpc.GetGroupResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.GetGroupResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.GetGroupRequest, ...grpc.CallOption) *charonrpc.GetGroupResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.GetGroupResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.GetGroupRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, in, opts
func (_m *GroupManagerClient) List(ctx context.Context, in *charonrpc.ListGroupsRequest, opts ...grpc.CallOption) (*charonrpc.ListGroupsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.ListGroupsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.ListGroupsRequest, ...grpc.CallOption) *charonrpc.ListGroupsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.ListGroupsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.ListGroupsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, in, opts
func (_m *GroupManagerClient) Delete(ctx context.Context, in *charonrpc.DeleteGroupRequest, opts ...grpc.CallOption) (*google_protobuf1.BoolValue, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *google_protobuf1.BoolValue
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.DeleteGroupRequest, ...grpc.CallOption) *google_protobuf1.BoolValue); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*google_protobuf1.BoolValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.DeleteGroupRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPermissions provides a mock function with given fields: ctx, in, opts
func (_m *GroupManagerClient) ListPermissions(ctx context.Context, in *charonrpc.ListGroupPermissionsRequest, opts ...grpc.CallOption) (*charonrpc.ListGroupPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.ListGroupPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.ListGroupPermissionsRequest, ...grpc.CallOption) *charonrpc.ListGroupPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.ListGroupPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.ListGroupPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPermissions provides a mock function with given fields: ctx, in, opts
func (_m *GroupManagerClient) SetPermissions(ctx context.Context, in *charonrpc.SetGroupPermissionsRequest, opts ...grpc.CallOption) (*charonrpc.SetGroupPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.SetGroupPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.SetGroupPermissionsRequest, ...grpc.CallOption) *charonrpc.SetGroupPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.SetGroupPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.SetGroupPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type PermissionManagerClient struct {
	mock.Mock
}

// Register provides a mock function with given fields: ctx, in, opts
func (_m *PermissionManagerClient) Register(ctx context.Context, in *charonrpc.RegisterPermissionsRequest, opts ...grpc.CallOption) (*charonrpc.RegisterPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.RegisterPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.RegisterPermissionsRequest, ...grpc.CallOption) *charonrpc.RegisterPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.RegisterPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.RegisterPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, in, opts
func (_m *PermissionManagerClient) List(ctx context.Context, in *charonrpc.ListPermissionsRequest, opts ...grpc.CallOption) (*charonrpc.ListPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.ListPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.ListPermissionsRequest, ...grpc.CallOption) *charonrpc.ListPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.ListPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.ListPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *PermissionManagerClient) Get(ctx context.Context, in *charonrpc.GetPermissionRequest, opts ...grpc.CallOption) (*charonrpc.GetPermissionResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.GetPermissionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.GetPermissionRequest, ...grpc.CallOption) *charonrpc.GetPermissionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.GetPermissionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.GetPermissionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type UserManagerClient struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in, opts
func (_m *UserManagerClient) Create(ctx context.Context, in *charonrpc.CreateUserRequest, opts ...grpc.CallOption) (*charonrpc.CreateUserResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.CreateUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.CreateUserRequest, ...grpc.CallOption) *charonrpc.CreateUserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.CreateUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.CreateUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Modify provides a mock function with given fields: ctx, in, opts
func (_m *UserManagerClient) Modify(ctx context.Context, in *charonrpc.ModifyUserRequest, opts ...grpc.CallOption) (*charonrpc.ModifyUserResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.ModifyUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.ModifyUserRequest, ...grpc.CallOption) *charonrpc.ModifyUserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.ModifyUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.ModifyUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *UserManagerClient) Get(ctx context.Context, in *charonrpc.GetUserRequest, opts ...grpc.CallOption) (*charonrpc.GetUserResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.GetUserResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.GetUserRequest, ...grpc.CallOption) *charonrpc.GetUserResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.GetUserResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.GetUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, in, opts
func (_m *UserManagerClient) List(ctx context.Context, in *charonrpc.ListUsersRequest, opts ...grpc.CallOption) (*charonrpc.ListUsersResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.ListUsersResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.ListUsersRequest, ...grpc.CallOption) *charonrpc.ListUsersResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.ListUsersResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.ListUsersRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, in, opts
func (_m *UserManagerClient) Delete(ctx context.Context, in *charonrpc.DeleteUserRequest, opts ...grpc.CallOption) (*google_protobuf1.BoolValue, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *google_protobuf1.BoolValue
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.DeleteUserRequest, ...grpc.CallOption) *google_protobuf1.BoolValue); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*google_protobuf1.BoolValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.DeleteUserRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPermissions provides a mock function with given fields: ctx, in, opts
func (_m *UserManagerClient) ListPermissions(ctx context.Context, in *charonrpc.ListUserPermissionsRequest, opts ...grpc.CallOption) (*charonrpc.ListUserPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.ListUserPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.ListUserPermissionsRequest, ...grpc.CallOption) *charonrpc.ListUserPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.ListUserPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.ListUserPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPermissions provides a mock function with given fields: ctx, in, opts
func (_m *UserManagerClient) SetPermissions(ctx context.Context, in *charonrpc.SetUserPermissionsRequest, opts ...grpc.CallOption) (*charonrpc.SetUserPermissionsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.SetUserPermissionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.SetUserPermissionsRequest, ...grpc.CallOption) *charonrpc.SetUserPermissionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.SetUserPermissionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.SetUserPermissionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGroups provides a mock function with given fields: ctx, in, opts
func (_m *UserManagerClient) ListGroups(ctx context.Context, in *charonrpc.ListUserGroupsRequest, opts ...grpc.CallOption) (*charonrpc.ListUserGroupsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.ListUserGroupsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.ListUserGroupsRequest, ...grpc.CallOption) *charonrpc.ListUserGroupsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.ListUserGroupsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.ListUserGroupsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetGroups provides a mock function with given fields: ctx, in, opts
func (_m *UserManagerClient) SetGroups(ctx context.Context, in *charonrpc.SetUserGroupsRequest, opts ...grpc.CallOption) (*charonrpc.SetUserGroupsResponse, error) {
	ret := _m.Called(ctx, in, opts)

	var r0 *charonrpc.SetUserGroupsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *charonrpc.SetUserGroupsRequest, ...grpc.CallOption) *charonrpc.SetUserGroupsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*charonrpc.SetUserGroupsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *charonrpc.SetUserGroupsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
