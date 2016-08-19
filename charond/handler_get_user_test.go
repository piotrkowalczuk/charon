package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/ntypes"
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
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsOwner,
				},
			},
			ent: userEntity{
				id:        2,
				createdBy: &ntypes.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
				},
			},
			ent: userEntity{
				id:        2,
				createdBy: &ntypes.Int64{Int64: 3, Valid: true},
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{
					id:          1,
					isSuperuser: true,
				},
			},
			ent: userEntity{
				id:          2,
				isSuperuser: true,
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{
					id:          1,
					isSuperuser: true,
				},
			},
			ent: userEntity{
				id: 2,
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{
					id: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: userEntity{
				id:        2,
				isStaff:   true,
				createdBy: &ntypes.Int64{Int64: 1, Valid: true},
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{
					id: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsStranger,
				},
			},
			ent: userEntity{
				id:        2,
				isStaff:   true,
				createdBy: &ntypes.Int64{Int64: 3, Valid: true},
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: userEntity{
				id: 1,
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{id: 1, isSuperuser: true},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: userEntity{
				id:          1,
				isSuperuser: true,
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
				user: &userEntity{id: 1},
			},
			ent: userEntity{
				id: 2,
			},
		},
		{
			req: charon.GetUserRequest{},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
			ent: userEntity{
				id:          2,
				isSuperuser: true,
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
