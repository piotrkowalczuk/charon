// +build unit,!postgres,!e2e

package main

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/nilt"
)

func TestGetUserHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.GetUserRequest
		act actor
		ent userEntity
	}{
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsOwner,
				},
			},
			ent: userEntity{
				ID:        2,
				CreatedBy: nilt.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
				},
			},
			ent: userEntity{
				ID:        2,
				CreatedBy: nilt.Int64{Int64: 3, Valid: true},
			},
		},
		{
			req: charon.GetUserRequest{},
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
			req: charon.GetUserRequest{},
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
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: userEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: nilt.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsStranger,
				},
			},
			ent: userEntity{
				ID:        2,
				IsStaff:   true,
				CreatedBy: nilt.Int64{Int64: 3, Valid: true},
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: userEntity{
				ID: 1,
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{ID: 1, IsSuperuser: true},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: userEntity{
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
		req charon.GetUserRequest
		act actor
		ent userEntity
	}{
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{},
			},
			ent: userEntity{},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
			ent: userEntity{
				ID: 2,
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: userEntity{
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
