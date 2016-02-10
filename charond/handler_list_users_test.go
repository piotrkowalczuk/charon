// +build unit !postgres

package main

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/nilt"
)

func TestListUsersHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.ListUsersRequest
		act actor
		ent userEntity
	}{
		{
			req: charon.ListUsersRequest{
				CreatedBy: &nilt.Int64{Int64: 1, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsOwner,
				},
			},
		},
		{
			req: charon.ListUsersRequest{
				CreatedBy: &nilt.Int64{Int64: 3, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
				},
			},
		},
		{
			req: charon.ListUsersRequest{
				IsSuperuser: &nilt.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
		},
		{
			req: charon.ListUsersRequest{},
			act: actor{
				user: &userEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
		},
		{
			req: charon.ListUsersRequest{
				IsStaff:   &nilt.Bool{Bool: true, Valid: true},
				CreatedBy: &nilt.Int64{Int64: 1, Valid: true},
			},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
		},
		{
			req: charon.ListUsersRequest{
				IsStaff:   &nilt.Bool{Bool: true, Valid: true},
				CreatedBy: &nilt.Int64{Int64: 3, Valid: true},
			},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsStranger,
				},
			},
		},
		{
			req: charon.ListUsersRequest{},
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
			req: charon.ListUsersRequest{
				IsSuperuser: &nilt.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 1, IsSuperuser: true},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
		},
	}

	h := &listUsersHandler{}
	for i, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error for %d: %s", i, err.Error())
		}
	}
}

func TestListUsersHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.ListUsersRequest
		act actor
		ent userEntity
	}{
		{
			req: charon.ListUsersRequest{},
			act: actor{
				user: &userEntity{},
			},
			ent: userEntity{},
		},
		{
			req: charon.ListUsersRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
			ent: userEntity{
				ID: 2,
			},
		},
		{
			req: charon.ListUsersRequest{
				IsSuperuser: &nilt.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
		},
	}

	h := &listUsersHandler{}
	for i, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error for %d, got nil", i)
		}
	}
}
