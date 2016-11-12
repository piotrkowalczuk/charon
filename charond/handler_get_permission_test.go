package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
)

func TestGetPermissionHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.GetPermissionRequest
		act actor
	}{
		{
			req: charonrpc.GetPermissionRequest{Id: 1},
			act: actor{
				user: &model.UserEntity{ID: 2},
				permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
				},
			},
		},
		{
			req: charonrpc.GetPermissionRequest{Id: 1},
			act: actor{
				user: &model.UserEntity{ID: 2, IsSuperuser: true},
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
		act actor
	}{
		{
			req: charonrpc.GetPermissionRequest{Id: 1},
			act: actor{
				user: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.GetPermissionRequest{Id: 1},
			act: actor{
				user: &model.UserEntity{ID: 2},
			},
		},
		{
			req: charonrpc.GetPermissionRequest{Id: 1},
			act: actor{
				user: &model.UserEntity{
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
