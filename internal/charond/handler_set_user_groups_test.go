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

func TestSetUserGroupsHandler_SetGroups(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	createGroupResp, err := suite.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{
		Name: "existing-group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	_, err = suite.charon.user.SetGroups(ctx, &charonrpc.SetUserGroupsRequest{
		Groups: []int64{createGroupResp.Group.Id},
		UserId: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestSetUserGroupsHandler_SetGroups_nonExistingGroup(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_, err := suite.charon.user.SetGroups(ctx, &charonrpc.SetUserGroupsRequest{
		Groups: []int64{1},
		UserId: 1,
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

func TestSetUserGroupsHandler_SetGroups_nonExistingUser(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	createGroupResp, err := suite.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{
		Name: "existing-group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	_, err = suite.charon.user.SetGroups(ctx, &charonrpc.SetUserGroupsRequest{
		Groups: []int64{createGroupResp.Group.Id},
		UserId: 2,
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

func TestSetUserGroupsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.SetUserGroupsRequest
		act session.Actor
	}{
		{
			req: charonrpc.SetUserGroupsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserGroupCanCreate,
					charon.UserGroupCanDelete,
				},
			},
		},
		{
			req: charonrpc.SetUserGroupsRequest{},

			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &setUserGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestSetUserGroupsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.SetUserGroupsRequest
		act session.Actor
	}{
		{
			req: charonrpc.SetUserGroupsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.SetUserGroupsRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.SetUserGroupsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserGroupCanDelete,
				},
			},
		},
		{
			req: charonrpc.SetUserGroupsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserGroupCanCreate,
				},
			},
		},
	}

	h := &setUserGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
