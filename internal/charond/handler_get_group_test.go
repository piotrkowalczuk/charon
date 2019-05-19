package charond

import (
	"database/sql"
	"testing"
	"time"

	"context"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/charon/internal/session/sessionmock"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetGroupHandler_Get_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cres := testRPCServerCreateGroup(t, suite, timeout(ctx), &charonrpc.CreateGroupRequest{
		Name:        "name",
		Description: ntypes.NewString("description"),
	})

	gres, err := suite.charon.group.Get(ctx, &charonrpc.GetGroupRequest{
		Id: cres.Group.Id,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if gres.Group.Name != cres.Group.Name {
		t.Errorf("wrong name, expected %s but got %s", cres.Group.Name, gres.Group.Name)
	}
	if gres.Group.Description != cres.Group.Description {
		t.Errorf("wrong description, expected %s but got %s", cres.Group.Description, gres.Group.Description)
	}
	_, err = suite.charon.group.Get(ctx, &charonrpc.GetGroupRequest{
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
}

func TestGetGroupHandler_Get_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	groupProviderMock := &modelmock.GroupProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.GetGroupRequest
		err  error
	}{
		"missing-group-id": {
			init: func(_ *testing.T) {},
			req:  charonrpc.GetGroupRequest{},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.GetGroupRequest{Id: 1},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"group-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				groupProviderMock.On("FindOneByID", mock.Anything, int64(1)).Return(nil, sql.ErrNoRows)
			},
			req: charonrpc.GetGroupRequest{Id: 1},
			err: grpcerr.E(codes.NotFound),
		},
		"group-fetch-canceled": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				groupProviderMock.On("FindOneByID", mock.Anything, int64(1)).Return(nil, context.Canceled)
			},
			req: charonrpc.GetGroupRequest{Id: 1},
			err: grpcerr.E(codes.Canceled),
		},
		"reverse-mapping-issue": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				groupProviderMock.On("FindOneByID", mock.Anything, int64(1)).Return(&model.GroupEntity{
					ID:        1,
					CreatedAt: time.Date(1, 1, 0, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			req: charonrpc.GetGroupRequest{Id: 1},
			err: grpcerr.E(codes.Internal),
		},
		"can-retrieve": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.GroupCanRetrieve},
						User:        &model.UserEntity{},
					}, nil).
					Once()
				groupProviderMock.On("FindOneByID", mock.Anything, int64(1)).Return(&model.GroupEntity{ID: 1}, nil)
			},
			req: charonrpc.GetGroupRequest{Id: 1},
		},
		"cannot-retrieve-without-permission": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsStaff: true}}, nil).
					Once()
			},
			req: charonrpc.GetGroupRequest{Id: 1},
			err: grpcerr.E(codes.PermissionDenied),
		},
	}

	h := getGroupHandler{
		handler: &handler{
			logger:        zap.L(),
			ActorProvider: actorProviderMock,
			repository: repositories{
				group: groupProviderMock,
			},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			actorProviderMock.ExpectedCalls = nil
			groupProviderMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.Get(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, actorProviderMock, groupProviderMock)
		})
	}
}
