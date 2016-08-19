package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/ntypes"
)

func TestCreateUserHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.CreateUserRequest
		act actor
	}{
		{
			req: charon.CreateUserRequest{},
			act: actor{
				user: &userEntity{id: 2},
				permissions: charon.Permissions{
					charon.UserCanCreate,
				},
			},
		},
		{
			req: charon.CreateUserRequest{
				IsStaff: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{id: 2},
				permissions: charon.Permissions{
					charon.UserCanCreateStaff,
				},
			},
		},
		{
			req: charon.CreateUserRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
			},
		},
		{
			req: charon.CreateUserRequest{},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
			},
		},
		{
			req: charon.CreateUserRequest{},
			act: actor{
				user:    &userEntity{},
				isLocal: true,
			},
		},
	}

	h := &createUserHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestCreateUserHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.CreateUserRequest
		act actor
	}{
		{
			req: charon.CreateUserRequest{},
			act: actor{
				user: &userEntity{},
			},
		},
		{
			req: charon.CreateUserRequest{},
			act: actor{
				user: &userEntity{id: 2},
			},
		},
		{
			req: charon.CreateUserRequest{},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
		{
			req: charon.CreateUserRequest{
				IsStaff: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
		{
			req: charon.CreateUserRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.UserCanCreateStaff,
				},
			},
		},
		{
			req: charon.CreateUserRequest{
				IsStaff: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{
					id: 2,
				},
				permissions: charon.Permissions{
					charon.UserCanCreate,
				},
			},
		},
	}

	h := &createUserHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
