package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/ntypes"
)

func TestDeleteUserHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.DeleteUserRequest
		act actor
		ent model.UserEntity
	}{
		{
			req: charonrpc.DeleteUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanDeleteAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				CreatedBy: ntypes.Int64{Int64: 3, Valid: true},
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
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
			req: charonrpc.DeleteUserRequest{},
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
			req: charonrpc.DeleteUserRequest{},
			act: actor{
				user: &model.UserEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: actor{
				user: &model.UserEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanDeleteStaffAsStranger,
				},
			},
			ent: model.UserEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: ntypes.Int64{Int64: 3, Valid: true},
			},
		},
	}

	h := &deleteUserHandler{}
	for i, d := range data {
		if err := h.firewall(&d.req, &d.act, &d.ent); err != nil {
			t.Errorf("unexpected error for %d: %s", i, err.Error())
		}
	}
}

func TestDeleteUserHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.DeleteUserRequest
		act actor
		ent model.UserEntity
	}{
		{
			req: charonrpc.DeleteUserRequest{},
			act: actor{
				user: &model.UserEntity{},
			},
			ent: model.UserEntity{},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
			},
			ent: model.UserEntity{
				ID: 2,
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
					charon.UserCanDeleteAsOwner,
					charon.UserCanDeleteStaffAsStranger,
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID: 1,
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
					charon.UserCanDeleteAsOwner,
					charon.UserCanDeleteStaffAsStranger,
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:          2,
				IsSuperuser: true,
			},
		},
		{
			req: charonrpc.DeleteUserRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1, IsSuperuser: true},
				permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
					charon.UserCanDeleteAsOwner,
					charon.UserCanDeleteStaffAsStranger,
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: model.UserEntity{
				ID:          1,
				IsSuperuser: true,
			},
		},
	}

	h := &deleteUserHandler{}
	for i, d := range data {
		if err := h.firewall(&d.req, &d.act, &d.ent); err == nil {
			t.Errorf("expected error for %d, got nil", i)
		}
	}
}
