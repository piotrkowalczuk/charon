package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
)

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
