package charond

import (
	"context"
	"testing"
	"time"

	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"

	"github.com/piotrkowalczuk/charon/internal/model"

	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/charon/internal/service"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynetest"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLoginHandler_Login_Unit(t *testing.T) {
	sessionMock := &mnemosynetest.SessionManagerClient{}
	userProviderMock := &modelmock.UserProvider{}
	refreshTokenProviderMock := &modelmock.RefreshTokenProvider{}
	hasher, err := password.NewBCryptHasher(5)
	if err != nil {
		t.Fatal(err)
	}
	pass, err := hasher.Hash([]byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.LoginRequest
		err  error
	}{
		"session-store-communication-failure": {
			init: func(t *testing.T) {
				usr := model.UserEntity{
					ID:          1,
					Username:    "test",
					Password:    pass,
					IsConfirmed: true,
					IsActive:    true,
				}
				userProviderMock.On("FindOneByUsername", mock.Anything, usr.Username).Return(&usr, nil)
				sessionMock.On("Start", mock.Anything, mock.Anything).Return(nil, status.Errorf(codes.Canceled, "example error"))
			},
			req: charonrpc.LoginRequest{Username: "test", Password: "test"},
			err: grpcerr.E(codes.Canceled),
		},
		"last-login-at-update-failure": {
			init: func(t *testing.T) {
				usr := model.UserEntity{
					ID:          1,
					Username:    "test",
					Password:    pass,
					IsConfirmed: true,
					IsActive:    true,
				}
				userProviderMock.On("FindOneByUsername", mock.Anything, usr.Username).Return(&usr, nil)
				sessionMock.On("Start", mock.Anything, mock.Anything).Return(&mnemosynerpc.StartResponse{
					Session: &mnemosynerpc.Session{
						AccessToken: "access_token",
					},
				}, nil)
				userProviderMock.On("UpdateLastLoginAt", mock.Anything, usr.ID).Return(int64(0), context.DeadlineExceeded)
			},
			req: charonrpc.LoginRequest{Username: "test", Password: "test"},
			err: grpcerr.E(codes.DeadlineExceeded),
		},
	}

	h := loginHandler{
		handler: &handler{
			logger:  zap.L(),
			session: sessionMock,
			repository: repositories{
				user: userProviderMock,
			},
		},
		userFinderFactory: &service.UserFinderFactory{
			UserRepository:         userProviderMock,
			RefreshTokenRepository: refreshTokenProviderMock,
			Hasher:                 hasher,
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			defer recoverTest(t)

			sessionMock.ExpectedCalls = nil
			userProviderMock.ExpectedCalls = nil
			refreshTokenProviderMock.ExpectedCalls = nil

			c.init(t)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := h.Login(ctx, &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t,
				sessionMock,
				userProviderMock,
				refreshTokenProviderMock,
			)
		})
	}
}

func TestLoginHandler_Login_E2E(t *testing.T) {
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
		"refresh_token_strategy": func(t *testing.T) {
			token, err := suite.charon.auth.Login(context.Background(), &charonrpc.LoginRequest{
				Strategy: &charonrpc.LoginRequest_RefreshToken{
					RefreshToken: &charonrpc.RefreshTokenStrategy{
						RefreshToken: "test",
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
