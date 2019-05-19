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
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

func TestModifyGroupHandler_Modify_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	res, err := suite.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{Name: "example"})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("missing-id", func(t *testing.T) {
		_, err := suite.charon.group.Modify(ctx, &charonrpc.ModifyGroupRequest{})
		assertErrorCode(t, err, codes.InvalidArgument, "group id is missing")
	})
	t.Run("nothing-to-modify", func(t *testing.T) {
		_, err := suite.charon.group.Modify(ctx, &charonrpc.ModifyGroupRequest{Id: res.Group.Id})
		assertErrorCode(t, err, codes.InvalidArgument, "nothing to be modified")
	})
	t.Run("ok", func(t *testing.T) {
		_, err := suite.charon.group.Modify(ctx, &charonrpc.ModifyGroupRequest{
			Id:          res.Group.Id,
			Name:        ntypes.NewString("A"),
			Description: ntypes.NewString("B"),
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestModifyGroupHandler_Modify_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	groupProviderMock := &modelmock.GroupProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.ModifyGroupRequest
		err  error
	}{
		"group-id-missing": {
			init: func(t *testing.T) {
			},
			req: charonrpc.ModifyGroupRequest{},
			err: grpcerr.E(codes.InvalidArgument),
		},
		"nothing-to-be-modified": {
			init: func(t *testing.T) {
			},
			req: charonrpc.ModifyGroupRequest{Id: 1},
			err: grpcerr.E(codes.InvalidArgument),
		},
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.ModifyGroupRequest{Id: 1, Name: ntypes.NewString("123")},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"storage-query-cancel": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				groupProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).
					Return(nil, context.Canceled).
					Once()
			},
			req: charonrpc.ModifyGroupRequest{Id: 1, Name: ntypes.NewString("123")},
			err: grpcerr.E(codes.Canceled),
		},
		"reverse-mapping-failure": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 2, IsSuperuser: true},
					}, nil).
					Once()
				groupProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).
					Return(&model.GroupEntity{
						ID:        1,
						Name:      "name",
						CreatedAt: brokenDate(),
					}, nil).
					Once()
			},
			req: charonrpc.ModifyGroupRequest{Id: 1, Name: ntypes.NewString("123")},
			err: grpcerr.E(codes.Internal),
		},
		"not-found": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 2, IsSuperuser: true},
					}, nil).
					Once()
				groupProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).
					Return(nil, sql.ErrNoRows).
					Once()
			},
			req: charonrpc.ModifyGroupRequest{Id: 1, Name: ntypes.NewString("123")},
			err: grpcerr.E(codes.NotFound),
		},
		"can-modify-as-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil).Once()
				groupProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).Return(&model.GroupEntity{
					ID:        1,
					Name:      "123",
					CreatedAt: time.Now(),
				}, nil).Once()
			},
			req: charonrpc.ModifyGroupRequest{Id: 1, Name: ntypes.NewString("123")},
		},
		"can-modify-with-permissions": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User:        &model.UserEntity{ID: 2},
					Permissions: charon.Permissions{charon.GroupCanModify},
				}, nil).Once()
				groupProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).Return(&model.GroupEntity{
					ID:        1,
					Name:      "123",
					CreatedAt: time.Now(),
				}, nil).Once()
			},
			req: charonrpc.ModifyGroupRequest{Id: 1, Name: ntypes.NewString("123")},
		},
		"cannot-modify-without-permissions": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2},
				}, nil).Once()
			},
			req: charonrpc.ModifyGroupRequest{Id: 1, Name: ntypes.NewString("123")},
			err: grpcerr.E(codes.PermissionDenied),
		},
	}

	h := modifyGroupHandler{
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
			defer recoverTest(t)

			actorProviderMock.ExpectedCalls = nil
			groupProviderMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.Modify(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, actorProviderMock, groupProviderMock)
		})
	}
}
