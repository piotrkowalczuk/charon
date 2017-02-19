package charond

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestBelongsToHandler_BelongsTo(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)
	resAct, err := suite.charon.auth.Actor(timeout(ctx), &wrappers.StringValue{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	var groups []int64
	for i := 0; i < 10; i++ {
		resGroup, err := suite.charon.group.Create(timeout(ctx), &charonrpc.CreateGroupRequest{
			Name: fmt.Sprintf("name-%d", i),
			Description: &ntypes.String{
				Valid: true,
				Chars: fmt.Sprintf("description-%d", i),
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		groups = append(groups, resGroup.Group.Id)
	}

	_, err = suite.charon.user.SetGroups(ctx, &charonrpc.SetUserGroupsRequest{
		UserId: resAct.Id,
		Groups: groups[:len(groups)/2],
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	cases := map[string]func(t *testing.T){
		"belongs": func(t *testing.T) {
			res, err := suite.charon.auth.BelongsTo(timeout(ctx), &charonrpc.BelongsToRequest{
				UserId:  resAct.Id,
				GroupId: groups[0],
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if !res.Value {
				t.Error("expected to belong")
			}
		},
		"not-belongs": func(t *testing.T) {
			res, err := suite.charon.auth.BelongsTo(timeout(ctx), &charonrpc.BelongsToRequest{
				UserId:  resAct.Id,
				GroupId: groups[len(groups)-1],
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.Value {
				t.Error("expected to not belong")
			}
		},
		"group-does-not-exists": func(t *testing.T) {
			res, err := suite.charon.auth.BelongsTo(timeout(ctx), &charonrpc.BelongsToRequest{
				UserId:  resAct.Id,
				GroupId: 99999999,
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.Value {
				t.Error("expected to not belong")
			}
		},
		"without-user-id": func(t *testing.T) {
			_, err := suite.charon.auth.BelongsTo(timeout(ctx), &charonrpc.BelongsToRequest{
				GroupId: groups[0],
			})
			if grpc.Code(err) != codes.InvalidArgument {
				t.Fatalf("wrong status code, expected %s but got %s", codes.InvalidArgument.String(), grpc.Code(err).String())
			}
		},

		"without-group-id": func(t *testing.T) {
			_, err := suite.charon.auth.BelongsTo(timeout(ctx), &charonrpc.BelongsToRequest{
				UserId: resAct.Id,
			})
			if grpc.Code(err) != codes.InvalidArgument {
				t.Fatalf("wrong status code, expected %s but got %s", codes.InvalidArgument.String(), grpc.Code(err).String())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

func TestBelongsToHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.BelongsToRequest
		act session.Actor
	}{
		{
			req: charonrpc.BelongsToRequest{UserId: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.UserGroupCanCheckBelongingAsStranger,
				},
			},
		},

		{
			req: charonrpc.BelongsToRequest{UserId: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},

		{
			req: charonrpc.BelongsToRequest{UserId: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
	}

	h := &belongsToHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestBelongsToHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.BelongsToRequest
		act session.Actor
	}{
		{
			req: charonrpc.BelongsToRequest{UserId: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
			},
		},
		{
			req: charonrpc.BelongsToRequest{UserId: 1},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &belongsToHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
