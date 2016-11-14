package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/ntypes"
)

func TestGetUserHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.GetUserRequest
		act actor
		ent model.UserEntity
	}{
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				CreatedBy: &ntypes.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				CreatedBy: &ntypes.Int64{Int64: 3, Valid: true},
			},
		},
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
			ent: model.UserEntity{
				ID:          2,
				IsSuperuser: true,
			},
		},
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
			ent: model.UserEntity{
				ID: 2,
			},
		},
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: &ntypes.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsStranger,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: &ntypes.Int64{Int64: 3, Valid: true},
			},
		},
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID: 1,
			},
		},
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1, IsSuperuser: true},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:          1,
				IsSuperuser: true,
			},
		},
	}

	h := &getUserHandler{}
	for i, d := range data {
		if err := h.firewall(&d.req, &d.act, &d.ent); err != nil {
			t.Errorf("unexpected error for %d: %s", i, err.Error())
		}
	}
}

func TestGetUserHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.GetUserRequest
		act actor
		ent model.UserEntity
	}{
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{},
			},
			ent: model.UserEntity{},
		},
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
			},
			ent: model.UserEntity{
				ID: 2,
			},
		},
		{
			req: charonrpc.GetUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:          2,
				IsSuperuser: true,
			},
		},
	}

	h := &getUserHandler{}
	for i, d := range data {
		if err := h.firewall(&d.req, &d.act, &d.ent); err == nil {
			t.Errorf("expected error for %d, got nil", i)
		}
	}
}
