package charond

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
)

func TestSetUserPermissionsHandler_SetPermissions(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)
	permissions := []string{
		"a:b:c",
	}

	_, err := suite.charon.permission.Register(ctx, &charonrpc.RegisterPermissionsRequest{
		Permissions: permissions,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = suite.charon.user.SetPermissions(ctx, &charonrpc.SetUserPermissionsRequest{
		UserId:      1,
		Permissions: permissions,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestSetUserPermissionsHandler_SetPermissions_nonExistingUser(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)
	permissions := []string{
		"a:b:c",
	}
	_, err := suite.charon.permission.Register(ctx, &charonrpc.RegisterPermissionsRequest{
		Permissions: permissions,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = suite.charon.user.SetPermissions(ctx, &charonrpc.SetUserPermissionsRequest{
		UserId:      2,
		Permissions: permissions,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() != codes.NotFound {
				t.Fatalf("wrong error code, expected %s but got %s for error: %s", codes.NotFound, st.Code(), err.Error())
			}
		} else {
			t.Errorf("wrong error type: %T", err)
		}
	}
}

func TestSetUserPermissionsHandler_SetPermissions_nonExistingPermission(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_, err := suite.charon.user.SetPermissions(ctx, &charonrpc.SetUserPermissionsRequest{
		UserId: 1,
		Permissions: []string{
			"fake:fake:fake",
		},
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() != codes.NotFound {
				t.Fatalf("wrong error code, expected %s but got %s for error: %s", codes.NotFound, st.Code(), err.Error())
			}
		} else {
			t.Errorf("wrong error type: %T", err)
		}
	}
	_, err = suite.charon.user.SetPermissions(ctx, &charonrpc.SetUserPermissionsRequest{
		UserId: 1,
		Permissions: []string{
			"fake:fake:fake",
		},
		Force: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestSetUserPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.SetUserPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserPermissionCanDelete,
					charon.UserPermissionCanCreate,
				},
			},
		},
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &setUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestSetUserPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.SetUserPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserPermissionCanDelete,
				},
			},
		},
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserPermissionCanCreate,
				},
			},
		},
	}

	h := &setUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
