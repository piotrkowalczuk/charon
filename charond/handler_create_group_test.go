package charond

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
)

func TestCreateGroupHandler_Create(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cases := map[string]func(t *testing.T){
		"full": func(t *testing.T) {
			req := &charonrpc.CreateGroupRequest{
				Name: "name-full",
				Description: &ntypes.String{
					Valid:  true,
					String: "description",
				},
			}
			res, err := suite.charon.group.Create(ctx, req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.Group.Name != req.Name {
				t.Errorf("wrong name, expected %s but got %s", req.Name, res.Group.Name)
			}
			if res.Group.Description != req.Description.StringOr("") {
				t.Errorf("wrong description, expected %#v but got %#v", req.Description.StringOr(""), res.Group.Description)
			}
		},
		"only-name": func(t *testing.T) {
			req := &charonrpc.CreateGroupRequest{
				Name: "name-only-name",
			}
			res, err := suite.charon.group.Create(ctx, req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.Group.Name != req.Name {
				t.Errorf("wrong name, expected %s but got %s", req.Name, res.Group.Name)
			}
			if res.Group.Description != "" {
				t.Errorf("wrong description, expected %#v but got %#v", "", res.Group.Description)
			}
		},
		"same-name-twice": func(t *testing.T) {
			req := &charonrpc.CreateGroupRequest{
				Name: "same-name-twice",
			}
			_, err := suite.charon.group.Create(ctx, req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			_, err = suite.charon.group.Create(ctx, req)
			if grpc.Code(err) != codes.AlreadyExists {
				t.Fatalf("wrong status code, expected %s but got %s", codes.AlreadyExists.String(), grpc.Code(err).String())
			}
		},
		"only-description": func(t *testing.T) {
			req := &charonrpc.CreateGroupRequest{
				Description: &ntypes.String{
					Valid:  true,
					String: "description",
				},
			}
			_, err := suite.charon.group.Create(ctx, req)
			if grpc.Code(err) != codes.InvalidArgument {
				t.Fatalf("wrong status code, expected %s but got %s", codes.InvalidArgument.String(), grpc.Code(err).String())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

func TestCreateGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.CreateGroupRequest
		act session.Actor
	}{
		{
			req: charonrpc.CreateGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.GroupCanCreate,
				},
			},
		},
		{
			req: charonrpc.CreateGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &createGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestCreateGroupHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.CreateGroupRequest
		act session.Actor
	}{
		{
			req: charonrpc.CreateGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
			},
		},
		{
			req: charonrpc.CreateGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.CreateGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
	}

	h := &createGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
