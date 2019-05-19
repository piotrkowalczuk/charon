package charond

import (
	"context"
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/charon/internal/session/sessionmock"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

func TestIsGrantedHandler_IsGranted_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_, err := suite.charon.user.SetPermissions(ctx, &charonrpc.SetUserPermissionsRequest{
		UserId:      1,
		Permissions: []string{charon.PermissionCanRetrieve.String()},
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := suite.charon.auth.IsGranted(ctx, &charonrpc.IsGrantedRequest{
		UserId:     1,
		Permission: charon.PermissionCanRetrieve.String(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.GetValue() {
		t.Error("should be granted")
	}
	res, err = suite.charon.auth.IsGranted(ctx, &charonrpc.IsGrantedRequest{
		UserId:     1,
		Permission: charon.PermissionCanDelete.String(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.GetValue() {
		t.Error("should not be granted")
	}
}

func TestIsGrantedHandler_IsGranted_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	userProviderMock := &modelmock.UserProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.IsGrantedRequest
		err  error
	}{
		"missing-user-id": {
			init: func(_ *testing.T) {},
			req:  charonrpc.IsGrantedRequest{Permission: "123"},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		"missing-permission": {
			init: func(_ *testing.T) {},
			req:  charonrpc.IsGrantedRequest{UserId: 1},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.IsGrantedRequest{UserId: 1, Permission: "123:123:123"},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"cannot-check-as-a-stranger-if-missing-permission": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{ID: 2, IsStaff: true}}, nil).
					Once()
			},
			req: charonrpc.IsGrantedRequest{UserId: 1, Permission: "123:123:123"},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-check-if-same-user-id": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{ID: 1}}, nil).
					Once()

				userProviderMock.On("IsGranted", mock.Anything, int64(1), charon.Permission("123:123:123")).
					Return(true, nil).
					Once()
			},
			req: charonrpc.IsGrantedRequest{UserId: 1, Permission: "123:123:123"},
		},
		"can-check-if-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{ID: 2, IsSuperuser: true}}, nil).
					Once()

				userProviderMock.On("IsGranted", mock.Anything, int64(1), charon.Permission("123:123:123")).
					Return(true, nil).
					Once()
			},
			req: charonrpc.IsGrantedRequest{UserId: 1, Permission: "123:123:123"},
		},
		"can-check-as-a-stranger-if-have-permission": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserPermissionCanCheckGrantingAsStranger},
						User:        &model.UserEntity{ID: 2},
					}, nil).
					Once()

				userProviderMock.On("IsGranted", mock.Anything, int64(1), charon.Permission("123:123:123")).
					Return(true, nil).
					Once()
			},
			req: charonrpc.IsGrantedRequest{UserId: 1, Permission: "123:123:123"},
		},
		"request-canceled": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{ID: 2, IsSuperuser: true}}, nil).
					Once()
				userProviderMock.On("IsGranted", mock.Anything, int64(1), charon.Permission("123:123:123")).
					Return(false, context.Canceled).
					Once()
			},
			req: charonrpc.IsGrantedRequest{UserId: 1, Permission: "123:123:123"},
			err: grpcerr.E(codes.Canceled),
		},
	}

	h := isGrantedHandler{
		handler: &handler{
			logger:        zap.L(),
			ActorProvider: actorProviderMock,
			repository: repositories{
				user: userProviderMock,
			},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			actorProviderMock.ExpectedCalls = nil
			userProviderMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.IsGranted(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, actorProviderMock, userProviderMock)
		})
	}
}
