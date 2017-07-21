package charond

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
)

func TestGetPermissionHandler_Get(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	gres, err := suite.charon.permission.Get(ctx, &charonrpc.GetPermissionRequest{
		Id: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if gres.Permission != charon.AllPermissions[0].String() {
		t.Errorf("wrong permission, expected %s but got %s", charon.AllPermissions[0], gres.Permission)
	}

	_, err = suite.charon.permission.Get(ctx, &charonrpc.GetPermissionRequest{
		Id: 1000,
	})
	if err == nil {
		t.Fatal("missing error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.NotFound {
			t.Errorf("wrong error code, expected %s but got %s", codes.NotFound.String(), st.Code().String())
		}
	}
	_, err = suite.charon.permission.Get(ctx, &charonrpc.GetPermissionRequest{
		Id: 0,
	})
	if err == nil {
		t.Fatal("missing error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.InvalidArgument {
			t.Errorf("wrong error code, expected %s but got %s", codes.InvalidArgument.String(), st.Code().String())
		}
	}
}

func TestGetPermissionHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.GetPermissionRequest
		act session.Actor
	}{
		{
			req: charonrpc.GetPermissionRequest{Id: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
				},
			},
		},
		{
			req: charonrpc.GetPermissionRequest{Id: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &getPermissionHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestGetPermissionHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.GetPermissionRequest
		act session.Actor
	}{
		{
			req: charonrpc.GetPermissionRequest{Id: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.GetPermissionRequest{Id: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
			},
		},
		{
			req: charonrpc.GetPermissionRequest{Id: 1},
			act: session.Actor{
				User: &model.UserEntity{
					ID: 2, IsStaff: true,
				},
			},
		},
	}

	h := &getPermissionHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
