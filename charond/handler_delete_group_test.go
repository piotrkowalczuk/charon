package charond

import (
	"fmt"
	"testing"

	"math"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestDeleteGroupHandler_Delete(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)
	resAct, err := suite.charon.auth.Actor(ctx, &wrappers.StringValue{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	var groups []int64
	for i := 0; i < 10; i++ {
		resGroup, err := suite.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{
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
		"not-assigned": func(t *testing.T) {
			done, err := suite.charon.group.Delete(ctx, &charonrpc.DeleteGroupRequest{
				Id: groups[len(groups)-1],
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if !done.Value {
				t.Error("group expected to be removed")
			}
		},
		"not-existing": func(t *testing.T) {
			_, err := suite.charon.group.Delete(ctx, &charonrpc.DeleteGroupRequest{
				Id: math.MaxInt64,
			})
			if grpc.Code(err) != codes.NotFound {
				t.Errorf("wrong status code, expected %s but got %s", codes.NotFound.String(), grpc.Code(err).String())
			}
		},
		"assigned": func(t *testing.T) {
			_, err := suite.charon.group.Delete(ctx, &charonrpc.DeleteGroupRequest{
				Id: groups[0],
			})
			if grpc.Code(err) != codes.FailedPrecondition {
				t.Errorf("wrong status code, expected %s but got %s", codes.FailedPrecondition.String(), grpc.Code(err).String())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

func TestDeleteGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.DeleteGroupRequest
		act session.Actor
	}{
		{
			req: charonrpc.DeleteGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.GroupCanDelete,
				},
			},
		},
		{
			req: charonrpc.DeleteGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &deleteGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestDeleteGroupHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.DeleteGroupRequest
		act session.Actor
	}{
		{
			req: charonrpc.DeleteGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
			},
		},
		{
			req: charonrpc.DeleteGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &deleteGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
