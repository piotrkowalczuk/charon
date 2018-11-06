package charond

import (
	"context"
	"database/sql"
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/charon/internal/session/sessionmock"
	"github.com/piotrkowalczuk/sklog"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetPermissionHandler_Get_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	gres, err := suite.charon.permission.Get(ctx, &charonrpc.GetPermissionRequest{
		Id: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if gres.Permission != charon.AllPermissions[0].String() {
		t.Errorf("wrong permission, expected %s but got %s", charon.AllPermissions[0], gres.Permission)
	}

	_, err = suite.charon.permission.Get(ctx, &charonrpc.GetPermissionRequest{
		Id: 1000,
	})
	if err == nil {
		t.Fatal("missing error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.NotFound {
			t.Errorf("wrong error code, expected %s but got %s", codes.NotFound.String(), st.Code().String())
		}
	}
	_, err = suite.charon.permission.Get(ctx, &charonrpc.GetPermissionRequest{
		Id: 0,
	})
	if err == nil {
		t.Fatal("missing error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.InvalidArgument {
			t.Errorf("wrong error code, expected %s but got %s", codes.InvalidArgument.String(), st.Code().String())
		}
	}
}

func TestGetPermissionHandler_Get_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	permissionProviderMock := &modelmock.PermissionProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.GetPermissionRequest
		err  error
	}{
		"missing-group-id": {
			init: func(_ *testing.T) {},
			req:  charonrpc.GetPermissionRequest{},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.GetPermissionRequest{Id: 1},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"permission-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				permissionProviderMock.On("FindOneByID", mock.Anything, int64(1)).Return(nil, sql.ErrNoRows)
			},
			req: charonrpc.GetPermissionRequest{Id: 1},
			err: grpcerr.E(codes.NotFound),
		},
		"permission-fetch-canceled": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				permissionProviderMock.On("FindOneByID", mock.Anything, int64(1)).Return(nil, context.Canceled)
			},
			req: charonrpc.GetPermissionRequest{Id: 1},
			err: grpcerr.E(codes.Canceled),
		},
		"can-retrieve": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.PermissionCanRetrieve},
						User:        &model.UserEntity{},
					}, nil).
					Once()
				permissionProviderMock.On("FindOneByID", mock.Anything, int64(1)).Return(&model.PermissionEntity{ID: 1}, nil)
			},
			req: charonrpc.GetPermissionRequest{Id: 1},
		},
		"cannot-retrieve-without-permission": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{
						ID:      1,
						IsStaff: true,
					}}, nil).
					Once()
			},
			req: charonrpc.GetPermissionRequest{Id: 1},
			err: grpcerr.E(codes.PermissionDenied),
		},
	}

	h := getPermissionHandler{
		handler: &handler{
			logger:        sklog.NewTestLogger(t),
			ActorProvider: actorProviderMock,
			repository: repositories{
				permission: permissionProviderMock,
			},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			actorProviderMock.ExpectedCalls = nil
			permissionProviderMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.Get(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, actorProviderMock, permissionProviderMock)
		})
	}
}
