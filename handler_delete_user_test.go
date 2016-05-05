package charon

import (
	"testing"

	"github.com/piotrkowalczuk/ntypes"
)

func TestDeleteUserHandler_firewall_success(t *testing.T) {
	data := []struct {
		req DeleteUserRequest
		act actor
		ent userEntity
	}{
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserCanDeleteAsOwner,
				},
			},
			ent: userEntity{
				ID:        2,
				CreatedBy: &ntypes.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserCanDeleteAsStranger,
				},
			},
			ent: userEntity{
				ID:        2,
				CreatedBy: &ntypes.Int64{Int64: 3, Valid: true},
			},
		},
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
			ent: userEntity{
				ID:          2,
				IsSuperuser: true,
			},
		},
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
			ent: userEntity{
				ID: 2,
			},
		},
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: Permissions{
					UserCanDeleteStaffAsOwner,
				},
			},
			ent: userEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: &ntypes.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: Permissions{
					UserCanDeleteStaffAsStranger,
				},
			},
			ent: userEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: &ntypes.Int64{Int64: 3, Valid: true},
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
		req DeleteUserRequest
		act actor
		ent userEntity
	}{
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{},
			},
			ent: userEntity{},
		},
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
			ent: userEntity{
				ID: 2,
			},
		},
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserCanDeleteAsStranger,
					UserCanDeleteAsOwner,
					UserCanDeleteStaffAsStranger,
					UserCanDeleteStaffAsOwner,
				},
			},
			ent: userEntity{
				ID: 1,
			},
		},
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserCanDeleteAsStranger,
					UserCanDeleteAsOwner,
					UserCanDeleteStaffAsStranger,
					UserCanDeleteStaffAsOwner,
				},
			},
			ent: userEntity{
				ID:          2,
				IsSuperuser: true,
			},
		},
		{
			req: DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1, IsSuperuser: true},
				permissions: Permissions{
					UserCanDeleteAsStranger,
					UserCanDeleteAsOwner,
					UserCanDeleteStaffAsStranger,
					UserCanDeleteStaffAsOwner,
				},
			},
			ent: userEntity{
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
