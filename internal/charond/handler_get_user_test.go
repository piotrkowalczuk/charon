package charond

import (
	"context"
	"testing"
	"time"

	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/charon/internal/session/sessionmock"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
)

func TestGetUserHandler_Get_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	res, err := suite.charon.user.Get(ctx, &charonrpc.GetUserRequest{Id: 1})
	if err != nil {
		t.Fatal(err)
	}
	if !res.GetUser().GetIsSuperuser() {
		t.Error("superuser should be returned")
	}
	if res.GetUser().GetUsername() != "test" {
		t.Error("username should match")
	}
}

func TestGetUserHandler_Get_Unit(t *testing.T) {
	userProviderMock := &modelmock.UserProvider{}
	actorProviderMock := &sessionmock.ActorProvider{}

	cases := []struct {
		req  charonrpc.GetUserRequest
		init func(*testing.T)
		err  error
	}{
		{
			req:  charonrpc.GetUserRequest{},
			init: func(_ *testing.T) {},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID:        2,
					CreatedBy: ntypes.Int64{Int64: 3, Valid: true},
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 1},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsStranger,
					},
				}, nil)
			},
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID:          2,
					IsSuperuser: true,
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{
						ID:          1,
						IsSuperuser: true,
					},
				}, nil)
			},
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID: 2,
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{
						ID:          1,
						IsSuperuser: true,
					},
				}, nil)
			},
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID:        2,
					IsStaff:   true,
					CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{
						ID: 1,
					},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveStaffAsOwner,
					},
				}, nil)
			},
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID:        2,
					IsStaff:   true,
					CreatedBy: ntypes.Int64{Int64: 3, Valid: true},
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{
						ID: 1,
					},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveStaffAsStranger,
					},
				}, nil)
			},
		},
		{
			req: charonrpc.GetUserRequest{Id: 1},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID: 1,
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 1},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsStranger,
						charon.UserCanRetrieveAsOwner,
						charon.UserCanRetrieveStaffAsStranger,
						charon.UserCanRetrieveStaffAsOwner,
					},
				}, nil)
			},
		},
		{
			req: charonrpc.GetUserRequest{Id: 1},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID:          1,
					IsSuperuser: true,
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 1, IsSuperuser: true},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsStranger,
						charon.UserCanRetrieveAsOwner,
						charon.UserCanRetrieveStaffAsStranger,
						charon.UserCanRetrieveStaffAsOwner,
					},
				}, nil)
			},
		},
		{
			req: charonrpc.GetUserRequest{Id: 1},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID:        1,
					CreatedBy: ntypes.Int64{Int64: 2, Valid: true},
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsOwner,
					},
				}, nil)
			},
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{},
				}, nil)
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{ID: 2}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 1},
				}, nil)
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID:      2,
					IsStaff: true,
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 1},
				}, nil)
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID:        2,
					CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
					IsStaff:   true,
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{
						ID: 1,
					},
				}, nil)
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID:          2,
					IsSuperuser: true,
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 1},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsStranger,
						charon.UserCanRetrieveAsOwner,
						charon.UserCanRetrieveStaffAsStranger,
						charon.UserCanRetrieveStaffAsOwner,
					},
				}, nil)
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(&model.UserEntity{
					ID:        2,
					CreatedAt: time.Date(1, 1, 0, 0, 0, 0, 0, time.UTC),
				}, nil)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{
						ID:          1,
						IsSuperuser: true,
					},
				}, nil)
			},
			err: grpcerr.E(codes.Internal),
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{
						ID:          1,
						IsSuperuser: true,
					},
				}, nil)
			},
			err: grpcerr.E(codes.NotFound),
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(nil, grpcerr.E(codes.Unauthenticated, "session not found"))
			},
			err: grpcerr.E(codes.Unauthenticated),
		},
		{
			req: charonrpc.GetUserRequest{Id: 2},
			init: func(_ *testing.T) {
				userProviderMock.On("FindOneByID", mock.Anything, mock.Anything).Return(nil, context.DeadlineExceeded)
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{
						ID:          1,
						IsSuperuser: true,
					},
				}, nil)
			},
			err: grpcerr.E(codes.DeadlineExceeded),
		},
	}

	h := &getUserHandler{
		handler: &handler{
			ActorProvider: actorProviderMock,
			repository: repositories{
				user: userProviderMock,
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			userProviderMock.ExpectedCalls = nil
			actorProviderMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.Get(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, userProviderMock, actorProviderMock)
		})
	}
}
