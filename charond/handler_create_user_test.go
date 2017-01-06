package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestCreateUserHandler_Create(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cases := map[string]func(t *testing.T){
		"full": func(t *testing.T) {
			req := &charonrpc.CreateUserRequest{
				Username:      "username-full",
				FirstName:     "first-name-full",
				LastName:      "last-name-full",
				PlainPassword: "plain-password-full",
			}
			res, err := suite.charon.user.Create(ctx, req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.User.Username != req.Username {
				t.Errorf("wrong username, expected %s but got %s", req.Username, res.User.Username)
			}
			if res.User.FirstName != req.FirstName {
				t.Errorf("wrong first name, expected %#v but got %#v", req.FirstName, res.User.FirstName)
			}
			if res.User.LastName != req.LastName {
				t.Errorf("wrong last name, expected %#v but got %#v", req.LastName, res.User.LastName)
			}
		},
		"same-username-twice": func(t *testing.T) {
			req := &charonrpc.CreateUserRequest{
				Username:      "username-same-username-twice",
				FirstName:     "first-name-same-username-twice",
				LastName:      "last-name-same-username-twice",
				PlainPassword: "plain-password-same-username-twice",
			}
			_, err := suite.charon.user.Create(ctx, req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			_, err = suite.charon.user.Create(ctx, req)
			if grpc.Code(err) != codes.AlreadyExists {
				t.Errorf("wrong status code, expected %s but got %s", codes.AlreadyExists.String(), grpc.Code(err).String())
			}
		},
		"another-superuser-without-actor": func(t *testing.T) {
			// TODO: change logic so local IP do not break it
			//req := &charonrpc.CreateUserRequest{
			//	Username:      "username-another-superuser",
			//	FirstName:     "first-name-another-superuser",
			//	LastName:      "last-name-another-superuser",
			//	PlainPassword: "plain-password-another-superuser",
			//	IsSuperuser:   ntypes.True(),
			//}
			//_, err := suite.charon.user.Create(context.Background(), req)
			//if grpc.Code(err) != codes.AlreadyExists {
			//	t.Errorf("wrong status code, expected %s but got %s", codes.AlreadyExists.String(), grpc.Code(err).String())
			//}
			//if grpc.ErrorDesc(err) != "initial superuser account already exists" {
			//	t.Errorf("wrong error message, expected '%s' but got '%s'", "initial superuser account already exists", grpc.ErrorDesc(err))
			//}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

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
