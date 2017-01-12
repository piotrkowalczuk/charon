package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
)

func TestSetUserPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.SetUserPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserPermissionCanDelete,
					charon.UserPermissionCanCreate,
				},
			},
		},
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &setUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestSetUserPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.SetUserPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserPermissionCanDelete,
				},
			},
		},
		{
			req: charonrpc.SetUserPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserPermissionCanCreate,
				},
			},
		},
	}

	h := &setUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
