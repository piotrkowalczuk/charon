// +build unit !postgres

package main

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/nilt"
)

func TestDeleteUserHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.DeleteUserRequest
		act actor
		ent userEntity
	}{
		{
			req: charon.DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanDeleteAsOwner,
				},
			},
			ent: userEntity{
				ID:        2,
				CreatedBy: nilt.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charon.DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
				},
			},
			ent: userEntity{
				ID:        2,
				CreatedBy: nilt.Int64{Int64: 3, Valid: true},
			},
		},
		{
			req: charon.DeleteUserRequest{},
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
			req: charon.DeleteUserRequest{},
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
			req: charon.DeleteUserRequest{},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: userEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: nilt.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charon.DeleteUserRequest{},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanDeleteStaffAsStranger,
				},
			},
			ent: userEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: nilt.Int64{Int64: 3, Valid: true},
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
		req charon.DeleteUserRequest
		act actor
		ent userEntity
	}{
		{
			req: charon.DeleteUserRequest{},
			act: actor{
				user: &userEntity{},
			},
			ent: userEntity{},
		},
		{
			req: charon.DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
			ent: userEntity{
				ID: 2,
			},
		},
		{
			req: charon.DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
					charon.UserCanDeleteAsOwner,
					charon.UserCanDeleteStaffAsStranger,
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: userEntity{
				ID: 1,
			},
		},
		{
			req: charon.DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
					charon.UserCanDeleteAsOwner,
					charon.UserCanDeleteStaffAsStranger,
					charon.UserCanDeleteStaffAsOwner,
				},
			},
			ent: userEntity{
				ID:          2,
				IsSuperuser: true,
			},
		},
		{
			req: charon.DeleteUserRequest{},
			act: actor{
				user: &userEntity{ID: 1, IsSuperuser: true},
				permissions: charon.Permissions{
					charon.UserCanDeleteAsStranger,
					charon.UserCanDeleteAsOwner,
					charon.UserCanDeleteStaffAsStranger,
					charon.UserCanDeleteStaffAsOwner,
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
