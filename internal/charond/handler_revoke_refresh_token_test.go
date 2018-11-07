package charond

import (
	"context"
	"testing"

	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/charon/internal/session/sessionmock"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynetest"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRevokeRefreshTokenHandler_Disable_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cases := map[string]func(t *testing.T){
		"simple": func(t *testing.T) {
			res, err := suite.charon.refreshToken.Create(timeout(ctx), &charonrpc.CreateRefreshTokenRequest{
				Notes: &ntypes.String{
					Valid: true,
					Chars: "note",
				},
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			res2, err := suite.charon.refreshToken.Revoke(timeout(ctx), &charonrpc.RevokeRefreshTokenRequest{
				Token:  res.RefreshToken.Token,
				UserId: res.RefreshToken.UserId,
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if !res2.RefreshToken.Revoked {
				t.Error("refresh token expected to be disabled")
			}
		},
		"missing-user-id": func(t *testing.T) {
			exp := "refresh token cannot be disabled, missing user id"
			res, err := suite.charon.refreshToken.Create(timeout(ctx), &charonrpc.CreateRefreshTokenRequest{})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			_, err = suite.charon.refreshToken.Revoke(timeout(ctx), &charonrpc.RevokeRefreshTokenRequest{
				Token: res.RefreshToken.Token,
			})
			if err == nil {
				t.Fatal("error expected")
			}

			got := status.Convert(err).Message()
			if got != exp {
				t.Errorf("wrong error, expected '%s' but got '%s'", exp, got)
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

func TestRevokeRefreshTokenHandler_Disable_Unit(t *testing.T) {
	sessionMock := &mnemosynetest.SessionManagerClient{}
	refreshTokenMock := &modelmock.RefreshTokenProvider{}
	actorProviderMock := &sessionmock.ActorProvider{}

	h := revokeRefreshTokenHandler{
		handler: &handler{
			logger:        zap.L(),
			ActorProvider: actorProviderMock,
			session:       sessionMock,
			repository: repositories{
				refreshToken: refreshTokenMock,
			},
		},
	}

	cases := map[string]struct {
		req  charonrpc.RevokeRefreshTokenRequest
		init func(*testing.T)
		err  error
	}{
		"missing-token": {
			init: func(t *testing.T) {},
			req:  charonrpc.RevokeRefreshTokenRequest{UserId: 1},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		"missing-user-id": {
			init: func(t *testing.T) {},
			req:  charonrpc.RevokeRefreshTokenRequest{Token: "123"},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"refresh-token-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(1)).Return(nil, sql.ErrNoRows)
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
			err: grpcerr.E(codes.NotFound),
		},
		"refresh-token-fetch-timeout": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(1)).Return(nil, context.DeadlineExceeded)
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
			err: grpcerr.E(codes.DeadlineExceeded),
		},
		"cannot-disable-as-a-stranger": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.RefreshTokenCanRevokeAsOwner},
						User:        &model.UserEntity{ID: 1},
					}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(2)).Return(&model.RefreshTokenEntity{
					UserID: 2,
				}, nil)
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 2, Token: "123"},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-disable-as-an-owner": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 1},
					}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(1)).Return(&model.RefreshTokenEntity{
					UserID: 1,
				}, nil)
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-disable-as-a-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 4, IsSuperuser: true},
					}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(1)).Return(&model.RefreshTokenEntity{
					UserID: 1,
				}, nil)
				refreshTokenMock.On("UpdateOneByToken", mock.Anything, "123", mock.Anything).Return(&model.RefreshTokenEntity{
					UserID:  1,
					Token:   "123",
					Revoked: true,
				}, nil)
				sessionMock.On("Delete", mock.Anything, mock.Anything).Return(&wrappers.Int64Value{Value: 1}, nil)
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
		},
		"can-disable-as-a-stranger": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.RefreshTokenCanRevokeAsStranger},
						User:        &model.UserEntity{ID: 4},
					}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(1)).Return(&model.RefreshTokenEntity{
					UserID: 1,
				}, nil)
				refreshTokenMock.On("UpdateOneByToken", mock.Anything, "123", mock.Anything).Return(&model.RefreshTokenEntity{
					UserID:  1,
					Token:   "123",
					Revoked: true,
				}, nil)
				sessionMock.On("Delete", mock.Anything, mock.Anything).Return(&wrappers.Int64Value{Value: 1}, nil)
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
		},
		"can-disable-as-an-owner": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.RefreshTokenCanRevokeAsOwner},
						User:        &model.UserEntity{ID: 1},
					}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(1)).Return(&model.RefreshTokenEntity{
					UserID: 1,
				}, nil)
				refreshTokenMock.On("UpdateOneByToken", mock.Anything, "123", mock.Anything).Return(&model.RefreshTokenEntity{
					UserID:  1,
					Token:   "123",
					Revoked: true,
				}, nil)
				sessionMock.On("Delete", mock.Anything, mock.Anything).Return(&wrappers.Int64Value{Value: 1}, nil)
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
		},
		"storage-returns-broken-entity": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.RefreshTokenCanRevokeAsOwner},
						User:        &model.UserEntity{ID: 1},
					}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(1)).Return(&model.RefreshTokenEntity{
					UserID: 1,
				}, nil)
				refreshTokenMock.On("UpdateOneByToken", mock.Anything, "123", mock.Anything).Return(&model.RefreshTokenEntity{
					UserID:    1,
					Token:     "123",
					Revoked:   true,
					CreatedAt: brokenDate(),
				}, nil)
				sessionMock.On("Delete", mock.Anything, mock.Anything).Return(&wrappers.Int64Value{Value: 1}, nil)
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
			err: grpcerr.E(codes.Internal),
		},
		"token-was-removed-after-during-execution": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 4, IsSuperuser: true},
					}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(1)).Return(&model.RefreshTokenEntity{
					UserID: 1,
				}, nil)
				refreshTokenMock.On("UpdateOneByToken", mock.Anything, "123", mock.Anything).Return(nil, sql.ErrNoRows)
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
			err: grpcerr.E(codes.NotFound),
		},
		"request-canceled-during-deletion": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 4, IsSuperuser: true},
					}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(1)).Return(&model.RefreshTokenEntity{
					UserID: 1,
				}, nil)
				refreshTokenMock.On("UpdateOneByToken", mock.Anything, "123", mock.Anything).Return(nil, context.Canceled)
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
			err: grpcerr.E(codes.Canceled),
		},
		"session-deletion-failure": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 4, IsSuperuser: true},
					}, nil).
					Once()
				refreshTokenMock.On("FindOneByTokenAndUserID", mock.Anything, "123", int64(1)).Return(&model.RefreshTokenEntity{
					UserID: 1,
				}, nil)
				refreshTokenMock.On("UpdateOneByToken", mock.Anything, "123", mock.Anything).Return(&model.RefreshTokenEntity{
					UserID:  1,
					Token:   "123",
					Revoked: true,
				}, nil)
				sessionMock.On("Delete", mock.Anything, mock.Anything).Return(nil, status.Errorf(codes.Aborted, "something went wrong"))
			},
			req: charonrpc.RevokeRefreshTokenRequest{UserId: 1, Token: "123"},
			err: grpcerr.E(codes.Aborted),
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			sessionMock.ExpectedCalls = nil
			refreshTokenMock.ExpectedCalls = nil
			actorProviderMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.Revoke(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t)
		})
	}

}
