package charond

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
)

func TestModifyGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.ModifyGroupRequest
		act session.Actor
	}{
		{
			req: charonrpc.ModifyGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.GroupCanModify,
				},
			},
		},
		{
			req: charonrpc.ModifyGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &modifyGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestModifyGroupHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.ModifyGroupRequest
		act session.Actor
	}{
		{
			req: charonrpc.ModifyGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
			},
		},
		{
			req: charonrpc.ModifyGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.ModifyGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
	}

	h := &modifyGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}

func TestModifyGroupHandler_Modify(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cres := testRPCServerCreateGroup(t, suite, timeout(ctx), &charonrpc.CreateGroupRequest{
		Name:        "NAME",
		Description: ntypes.NewString("DESCRIPTION"),
	})
	req := &charonrpc.ModifyGroupRequest{
		Id:          cres.Group.Id,
		Name:        ntypes.NewString("name"),
		Description: ntypes.NewString("description"),
	}
	mres, err := suite.charon.group.Modify(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if mres.Group.Name != req.Name.StringOr("") {
		t.Errorf("wrong name, expected %s but got %s", req.Name.StringOr(""), mres.Group.Name)
	}
	if mres.Group.Description != req.Description.StringOr("") {
		t.Errorf("wrong description, expected %s but got %s", req.Description.StringOr(""), mres.Group.Description)
	}
	_, err = suite.charon.group.Modify(ctx, &charonrpc.ModifyGroupRequest{
		Id:          1000,
		Name:        ntypes.NewString("name"),
		Description: ntypes.NewString("description"),
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
