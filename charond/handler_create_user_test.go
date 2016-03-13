// +build unit,!postgres,!e2e

package main

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/nilt"
)

func TestCreateUserHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.CreateUserRequest
		act actor
	}{
		{
			req: charon.CreateUserRequest{},
			act: actor{
				user: &userEntity{ID: 2},
				permissions: charon.Permissions{
					charon.UserCanCreate,
				},
			},
		},
		{
			req: charon.CreateUserRequest{
				IsStaff: &nilt.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 2},
				permissions: charon.Permissions{
					charon.UserCanCreateStaff,
				},
			},
		},
		{
			req: charon.CreateUserRequest{
				IsSuperuser: &nilt.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
			},
		},
		{
			req: charon.CreateUserRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
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
				user: &userEntity{ID: 2},
			},
		},
		{
			req: charon.CreateUserRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charon.CreateUserRequest{
				IsStaff: &nilt.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charon.CreateUserRequest{
				IsSuperuser: &nilt.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanCreateStaff,
				},
			},
		},
		{
			req: charon.CreateUserRequest{
				IsStaff: &nilt.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{
					ID: 2,
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
			t.Errorf("expected error, got nil")
		}
	}
}
