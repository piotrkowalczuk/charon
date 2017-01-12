package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
)

func TestListGroupPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.ListGroupPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.ListGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.GroupPermissionCanRetrieve,
				},
			},
		},
		{
			req: charonrpc.ListGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &listGroupPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestListGroupPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.ListGroupPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.ListGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.ListGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.ListGroupPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{},
				Permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
					charon.UserPermissionCanRetrieve,
				},
			},
		},
	}

	h := &listGroupPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
