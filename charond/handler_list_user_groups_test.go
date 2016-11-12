package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
)

func TestListUserGroupsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.ListUserGroupsRequest
		act actor
	}{
		{
			req: charonrpc.ListUserGroupsRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserGroupCanRetrieve,
				},
			},
		},
		{
			req: charonrpc.ListUserGroupsRequest{},
			act: actor{
				user: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &listUserGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestListUserGroupsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.ListUserGroupsRequest
		act actor
	}{
		{
			req: charonrpc.ListUserGroupsRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.ListUserGroupsRequest{},
			act: actor{
				user: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.ListUserGroupsRequest{},
			act: actor{
				user: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
				permissions: charon.Permissions{
					charon.GroupCanRetrieve,
				},
			},
		},
	}

	h := &listUserGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
