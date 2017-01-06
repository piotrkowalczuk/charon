package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
)

func TestCreateUserHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.CreateUserRequest
		act session.Actor
	}{
		{
			req: charonrpc.CreateUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.UserCanCreate,
				},
			},
		},
		{

			req: charonrpc.CreateUserRequest{
				IsStaff: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.UserCanCreateStaff,
				},
			},
		},
		{
			req: charonrpc.CreateUserRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
		{
			req: charonrpc.CreateUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
		{
			req: charonrpc.CreateUserRequest{},
			act: session.Actor{
				User:    &model.UserEntity{},
				IsLocal: true,
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
		req charonrpc.CreateUserRequest
		act session.Actor
	}{
		{
			req: charonrpc.CreateUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{},
			},
		},
		{
			req: charonrpc.CreateUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
			},
		},
		{
			req: charonrpc.CreateUserRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.CreateUserRequest{
				IsStaff: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.CreateUserRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.UserCanCreateStaff,
				},
			},
		},
		{
			req: charonrpc.CreateUserRequest{
				IsStaff: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: session.Actor{
				User: &model.UserEntity{
					ID: 2,
				},
				Permissions: charon.Permissions{
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
