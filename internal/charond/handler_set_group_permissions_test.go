package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSetGroupPermissionsHandler_SetPermissions(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)
	permissions := []string{
		"a:b:c",
	}

	createGroupResp, err := suite.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{
		Name: "existing-group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = suite.charon.permission.Register(ctx, &charonrpc.RegisterPermissionsRequest{
		Permissions: permissions,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = suite.charon.group.SetPermissions(ctx, &charonrpc.SetGroupPermissionsRequest{
		GroupId:     createGroupResp.Group.Id,
		Permissions: permissions,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestSetGroupPermissionsHandler_SetPermissions_nonExistingGroup(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_, err := suite.charon.group.SetPermissions(ctx, &charonrpc.SetGroupPermissionsRequest{
		GroupId: 1,
		Permissions: []string{
			"fake",
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
}

func TestSetGroupPermissionsHandler_SetPermissions_nonExistingPermission(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	res, err := suite.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{
		Name: "existing-group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = suite.charon.group.SetPermissions(ctx, &charonrpc.SetGroupPermissionsRequest{
		GroupId: res.Group.Id,
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

	_, err = suite.charon.group.SetPermissions(ctx, &charonrpc.SetGroupPermissionsRequest{
		GroupId: res.Group.Id,
		Permissions: []string{
			"fake:fake:fake",
		},
		Force: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestSetGroupPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.SetGroupPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.SetGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.GroupPermissionCanCreate,
					charon.GroupPermissionCanDelete,
				},
			},
		},
		{
			req: charonrpc.SetGroupPermissionsRequest{},

			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &setGroupPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestSetGroupPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.SetGroupPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.SetGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.SetGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.SetGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.GroupPermissionCanDelete,
				},
			},
		},
		{
			req: charonrpc.SetGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.GroupPermissionCanCreate,
				},
			},
		},
	}

	h := &setGroupPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
