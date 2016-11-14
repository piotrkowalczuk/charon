package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
)

func TestListUserPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.ListUserPermissionsRequest
		act actor
	}{
		{
			req: charonrpc.ListUserPermissionsRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},

				permissions: charon.Permissions{
					charon.UserPermissionCanRetrieve,
				},
			},
		},
		{
			req: charonrpc.ListUserPermissionsRequest{},
			act: actor{
				user: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &listUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestListUserPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.ListUserPermissionsRequest
		act actor
	}{
		{
			req: charonrpc.ListUserPermissionsRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.ListUserPermissionsRequest{},
			act: actor{
				user: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.ListUserPermissionsRequest{},
			act: actor{
				user: &model.UserEntity{},
				permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
					charon.GroupPermissionCanRetrieve,
				},
			},
		},
	}

	h := &listUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
