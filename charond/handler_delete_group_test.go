package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
)

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
