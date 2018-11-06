package charond

import (
	"context"
	"testing"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/ntypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLoginHandler_Login(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cases := map[string]func(t *testing.T){
		"without-username": func(t *testing.T) {
			_, err := suite.charon.auth.Login(context.Background(), &charonrpc.LoginRequest{Password: "test"})
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("wrong status code, expected %s but got %s", codes.InvalidArgument.String(), status.Code(err).String())
			}
		},
		"without-password": func(t *testing.T) {
			_, err := suite.charon.auth.Login(context.Background(), &charonrpc.LoginRequest{Username: "test"})
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("wrong status code, expected %s but got %s", codes.InvalidArgument.String(), status.Code(err).String())
			}
		},
		"username_and_password_deprecated": func(t *testing.T) {
			token, err := suite.charon.auth.Login(context.Background(), &charonrpc.LoginRequest{Username: "test", Password: "test"})
			if err != nil {
				t.Fatalf("unexpected error: %s: with code %s", status.Convert(err).Message(), status.Code(err))
			}
			if len(token.Value) == 0 {
				t.Error("token should not be empty")
			}
		},
		"username_and_password_strategy": func(t *testing.T) {
			token, err := suite.charon.auth.Login(context.Background(), &charonrpc.LoginRequest{
				Strategy: &charonrpc.LoginRequest_UsernameAndPassword{
					UsernameAndPassword: &charonrpc.UsernameAndPasswordStrategy{
						Username: "test", Password: "test",
					},
				},
			})
			if err != nil {
				t.Fatalf("unexpected error: %s: with code %s", status.Convert(err).Message(), status.Code(err))
			}
			if len(token.Value) == 0 {
				t.Error("token should not be empty")
			}
		},
		"does-not-exists": func(t *testing.T) {
			_, err := suite.charon.auth.Login(context.Background(), &charonrpc.LoginRequest{Username: "test-not-exists", Password: "test"})
			if status.Code(err) != codes.Unauthenticated {
				t.Fatalf("wrong status code, expected %s but got %s", codes.Unauthenticated.String(), status.Code(err).String())
			}
		},
		"wrong-password": func(t *testing.T) {
			_, err := suite.charon.auth.Login(context.Background(), &charonrpc.LoginRequest{Username: "test", Password: "wrong-password"})
			if status.Code(err) != codes.Unauthenticated {
				t.Fatalf("wrong status code, expected %s but got %s", codes.Unauthenticated.String(), status.Code(err).String())
			}
		},
		"not-confirmed": func(t *testing.T) {
			req := &charonrpc.CreateUserRequest{
				Username:      "username-not-confirmed",
				FirstName:     "first-name-not-confirmed",
				LastName:      "last-name-not-confirmed",
				PlainPassword: "plain-password-not-confirmed",
				IsActive:      &ntypes.Bool{Bool: true, Valid: true},
			}
			_, err := suite.charon.user.Create(ctx, req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			_, err = suite.charon.auth.Login(context.Background(), &charonrpc.LoginRequest{Username: req.Username, Password: req.PlainPassword})
			if status.Code(err) != codes.Unauthenticated {
				t.Fatalf("wrong status code, expected %s but got %s", codes.Unauthenticated.String(), status.Code(err).String())
			}
		},
		"not-active": func(t *testing.T) {
			req := &charonrpc.CreateUserRequest{
				Username:      "username-not-active",
				FirstName:     "first-name-not-active",
				LastName:      "last-name-not-active",
				PlainPassword: "plain-password-not-active",
				IsConfirmed:   &ntypes.Bool{Bool: true, Valid: true},
			}
			_, err := suite.charon.user.Create(ctx, req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			_, err = suite.charon.auth.Login(context.Background(), &charonrpc.LoginRequest{Username: req.Username, Password: req.PlainPassword})
			if status.Code(err) != codes.Unauthenticated {
				t.Fatalf("wrong status code, expected %s but got %s", codes.Unauthenticated.String(), status.Code(err).String())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}
