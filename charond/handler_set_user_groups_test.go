package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
)

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
