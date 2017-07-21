package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetGroupHandler_Get(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cres := testRPCServerCreateGroup(t, suite, timeout(ctx), &charonrpc.CreateGroupRequest{
		Name:        "name",
		Description: ntypes.NewString("description"),
	})

	gres, err := suite.charon.group.Get(ctx, &charonrpc.GetGroupRequest{
		Id: cres.Group.Id,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if gres.Group.Name != cres.Group.Name {
		t.Errorf("wrong name, expected %s but got %s", cres.Group.Name, gres.Group.Name)
	}
	if gres.Group.Description != cres.Group.Description {
		t.Errorf("wrong description, expected %s but got %s", cres.Group.Description, gres.Group.Description)
	}
	_, err = suite.charon.group.Get(ctx, &charonrpc.GetGroupRequest{
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
}

func TestGetGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.GetGroupRequest
		act session.Actor
	}{
		{
			req: charonrpc.GetGroupRequest{Id: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.GroupCanRetrieve,
				},
			},
		},
		{
			req: charonrpc.GetGroupRequest{Id: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &getGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestGetGroupHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.GetGroupRequest
		act session.Actor
	}{
		{
			req: charonrpc.GetGroupRequest{Id: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.GetGroupRequest{Id: 1},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &getGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
