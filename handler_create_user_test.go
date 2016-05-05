package charon

import (
	"testing"

	"github.com/piotrkowalczuk/ntypes"
)

func TestCreateUserHandler_firewall_success(t *testing.T) {
	data := []struct {
		req CreateUserRequest
		act actor
	}{
		{
			req: CreateUserRequest{},
			act: actor{
				user: &userEntity{ID: 2},
				permissions: Permissions{
					UserCanCreate,
				},
			},
		},
		{
			req: CreateUserRequest{
				IsStaff: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 2},
				permissions: Permissions{
					UserCanCreateStaff,
				},
			},
		},
		{
			req: CreateUserRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
			},
		},
		{
			req: CreateUserRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
			},
		},
		{
			req: CreateUserRequest{},
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
		req CreateUserRequest
		act actor
	}{
		{
			req: CreateUserRequest{},
			act: actor{
				user: &userEntity{},
			},
		},
		{
			req: CreateUserRequest{},
			act: actor{
				user: &userEntity{ID: 2},
			},
		},
		{
			req: CreateUserRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: CreateUserRequest{
				IsStaff: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: CreateUserRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserCanCreateStaff,
				},
			},
		},
		{
			req: CreateUserRequest{
				IsStaff: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{
					ID: 2,
				},
				permissions: Permissions{
					UserCanCreate,
				},
			},
		},
	}

	h := &createUserHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
