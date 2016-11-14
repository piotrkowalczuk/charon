package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
)

func TestGetGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.GetGroupRequest
		act actor
	}{
		{
			req: charonrpc.GetGroupRequest{Id: 1},
			act: actor{
				user: &model.UserEntity{ID: 2},
				permissions: charon.Permissions{
					charon.GroupCanRetrieve,
				},
			},
		},
		{
			req: charonrpc.GetGroupRequest{Id: 1},
			act: actor{
				user: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &getGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestGetGroupHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.GetGroupRequest
		act actor
	}{
		{
			req: charonrpc.GetGroupRequest{Id: 1},
			act: actor{
				user: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.GetGroupRequest{Id: 1},
			act: actor{
				user: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &getGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
