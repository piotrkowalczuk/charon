package charond

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/charon/internal/session/sessionmock"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
)

func TestListRefreshTokensHandler_List_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_, err := suite.charon.refreshToken.Create(ctx, &charonrpc.CreateRefreshTokenRequest{
		ExpireAt: ptypes.TimestampNow(),
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := suite.charon.refreshToken.List(ctx, &charonrpc.ListRefreshTokensRequest{})
	if err != nil {
		t.Fatal(err)
	}

	// If daemon runs in test mode, it creates one test refresh token.
	if len(res.RefreshTokens) != 2 {
		t.Errorf("wrong number of refresh tokens, expected 1 got %d", len(res.RefreshTokens))
	}
}

func TestListRefreshTokensHandler_List_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	refreshTokenMock := &modelmock.RefreshTokenProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.ListRefreshTokensRequest
		err  error
	}{
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.ListRefreshTokensRequest{},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"reverse-mapping-failure": {
			init: func(*testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User:        &model.UserEntity{ID: 1},
					Permissions: charon.Permissions{charon.RefreshTokenCanRetrieveAsStranger},
				}, nil)
				refreshTokenMock.On("Find", mock.Anything, mock.Anything).Return([]*model.RefreshTokenEntity{
					{
						UserID:    1,
						Token:     "abc",
						CreatedAt: brokenDate(),
					},
				}, nil)
			},
			req: charonrpc.ListRefreshTokensRequest{},
			err: grpcerr.E(codes.Internal),
		},
		"context-canceled": {
			init: func(*testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User:        &model.UserEntity{ID: 1},
					Permissions: charon.Permissions{charon.RefreshTokenCanRetrieveAsStranger},
				}, nil)
				refreshTokenMock.On("Find", mock.Anything, mock.Anything).Return(nil, context.Canceled)
			},
			req: charonrpc.ListRefreshTokensRequest{},
			err: grpcerr.E(codes.Canceled),
		},
		"as-stranger": {
			init: func(*testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User:        &model.UserEntity{ID: 1},
					Permissions: charon.Permissions{charon.RefreshTokenCanRetrieveAsStranger},
				}, nil)
				refreshTokenMock.On("Find", mock.Anything, mock.Anything).Return([]*model.RefreshTokenEntity{
					{
						UserID: 1,
						Token:  "abc",
					},
				}, nil)
			},
			req: charonrpc.ListRefreshTokensRequest{},
		},
		"as-superuser": {
			init: func(*testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				refreshTokenMock.On("Find", mock.Anything, mock.Anything).Return([]*model.RefreshTokenEntity{
					{
						UserID: 1,
						Token:  "abc",
					},
				}, nil)
			},
			req: charonrpc.ListRefreshTokensRequest{},
		},
		"as-owner": {
			init: func(*testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User:        &model.UserEntity{ID: 2},
					Permissions: charon.Permissions{charon.RefreshTokenCanRetrieveAsOwner},
				}, nil)
				refreshTokenMock.On("Find", mock.Anything, mock.Anything).Return([]*model.RefreshTokenEntity{
					{
						UserID: 1,
						Token:  "abc",
					},
				}, nil)
			},
			req: charonrpc.ListRefreshTokensRequest{},
		},
		"as-no-one": {
			init: func(*testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 1},
				}, nil)
			},
			req: charonrpc.ListRefreshTokensRequest{},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"as-staff": {
			init: func(*testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{
						ID:      2,
						IsStaff: true,
					},
				}, nil)
			},
			req: charonrpc.ListRefreshTokensRequest{},
			err: grpcerr.E(codes.PermissionDenied),
		},
	}

	h := listRefreshTokensHandler{
		handler: &handler{
			ActorProvider: actorProviderMock,
			repository: repositories{
				refreshToken: refreshTokenMock,
			},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			actorProviderMock.ExpectedCalls = nil
			refreshTokenMock.ExpectedCalls = nil

			c.init(t)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := h.List(ctx, &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t,
				actorProviderMock,
				refreshTokenMock,
			)
		})
	}
}
