package charond

import (
	"context"
	"fmt"
	"testing"

	"net"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/charon/internal/session/sessionmock"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/qtypes"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
)

func TestListUsersHandler_List_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	t.Run("one", func(t *testing.T) {
		res, err := suite.charon.user.List(ctx, &charonrpc.ListUsersRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Users) != 1 {
			t.Errorf("wrong number of entities, expected %d but got %d", 1, len(res.Users))
		}
	})
	t.Run("unauthenticated", func(t *testing.T) {
		ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: &net.IPAddr{}})
		_, err := suite.charon.user.List(ctx, &charonrpc.ListUsersRequest{})
		assertErrorCode(t, err, codes.Unauthenticated, "mnemosyned: missing access token in metadata")
	})

	max := 10

	for i := 0; i < max; i++ {
		suite.charon.user.Create(ctx, &charonrpc.CreateUserRequest{
			Username:      fmt.Sprintf("username-%d", i),
			PlainPassword: fmt.Sprintf("password-%d", i),
			FirstName:     fmt.Sprintf("first-name-%d", i),
			LastName:      fmt.Sprintf("last-name-%d", i),
		})
	}

	t.Run("offset", func(t *testing.T) {
		// total number of 11 users (superuser + 10 auto generated)
		res, err := suite.charon.user.List(ctx, &charonrpc.ListUsersRequest{
			Offset: ntypes.NewInt64(int64(max)),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Users) != 1 {
			t.Errorf("wrong number of entities, expected %d but got %d:\n%v", 1, len(res.Users), res.Users)
		}
	})
	t.Run("order-by", func(t *testing.T) {
		res, err := suite.charon.user.List(ctx, &charonrpc.ListUsersRequest{
			IsSuperuser: ntypes.False(),
			OrderBy: []*charonrpc.Order{
				{
					Name:       model.TableUserColumnFirstName,
					Descending: true,
				},
				{
					Name:       model.TableUserColumnLastName,
					Descending: true,
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		exp := "username-9"
		if res.Users[0].Username != exp {
			t.Errorf("wrong group name, expected %s but got %s", exp, res.Users[0])
		}
	})
}

func TestListUsersHandler_List_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	userProviderMock := &modelmock.UserProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.ListUsersRequest
		err  error
	}{
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.ListUsersRequest{},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"storage-query-cancel": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				userProviderMock.On("Find", mock.Anything, mock.Anything).Return(nil, context.Canceled)
			},
			req: charonrpc.ListUsersRequest{},
			err: grpcerr.E(codes.Canceled),
		},
		"reverse-mapping-failure": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				userProviderMock.On("Find", mock.Anything, mock.Anything).Return([]*model.UserEntity{
					{
						ID:        1,
						LastName:  "last-name",
						FirstName: "first-name",
						CreatedAt: brokenDate(),
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{},
			err: grpcerr.E(codes.Internal),
		},
		"can-retrieve-superuser-as-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				userProviderMock.On("Find", mock.Anything, mock.Anything).Return([]*model.UserEntity{
					{
						ID:        1,
						LastName:  "last-name",
						FirstName: "first-name",
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{IsSuperuser: ntypes.True()},
		},
		"cannot-retrieve-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsStranger,
						charon.UserCanRetrieveAsOwner,
						charon.UserCanRetrieveStaffAsStranger,
						charon.UserCanRetrieveStaffAsOwner,
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{IsSuperuser: ntypes.True()},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-retrieve-as-owner": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsOwner,
					},
				}, nil)
				userProviderMock.On("Find", mock.Anything, mock.Anything).Return([]*model.UserEntity{
					{
						ID:        1,
						CreatedBy: ntypes.Int64{Int64: 2, Valid: true},
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{CreatedBy: qtypes.EqualInt64(2)},
		},
		"cannot-retrieve-as-owner": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveStaffAsOwner,
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{CreatedBy: qtypes.EqualInt64(2)},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-retrieve-staff-as-owner": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveStaffAsOwner,
					},
				}, nil)
				userProviderMock.On("Find", mock.Anything, mock.Anything).Return([]*model.UserEntity{
					{
						ID:        1,
						IsStaff:   true,
						CreatedBy: ntypes.Int64{Int64: 2, Valid: true},
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{CreatedBy: qtypes.EqualInt64(2), IsStaff: ntypes.True()},
		},
		"cannot-retrieve-staff-as-owner": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsOwner,
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{CreatedBy: qtypes.EqualInt64(2), IsStaff: ntypes.True()},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-retrieve-staff-as-stranger": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 3},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveStaffAsStranger,
					},
				}, nil)
				userProviderMock.On("Find", mock.Anything, mock.Anything).Return([]*model.UserEntity{
					{
						ID:        1,
						IsStaff:   true,
						CreatedBy: ntypes.Int64{Int64: 2, Valid: true},
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{CreatedBy: qtypes.EqualInt64(2), IsStaff: ntypes.True()},
		},
		"cannot-retrieve-staff-as-stranger": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 3},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveStaffAsOwner,
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{CreatedBy: qtypes.EqualInt64(2), IsStaff: ntypes.True()},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-retrieve-as-stranger-search-by-creator": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 3},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsStranger,
					},
				}, nil)
				userProviderMock.On("Find", mock.Anything, mock.Anything).Return([]*model.UserEntity{
					{
						ID:        1,
						CreatedBy: ntypes.Int64{Int64: 2, Valid: true},
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{CreatedBy: qtypes.EqualInt64(2)},
		},
		"can-retrieve-as-stranger": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 3},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsStranger,
					},
				}, nil)
				userProviderMock.On("Find", mock.Anything, mock.Anything).Return([]*model.UserEntity{
					{
						ID:        1,
						CreatedBy: ntypes.Int64{Int64: 2, Valid: true},
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{},
		},
		"cannot-retrieve-as-stranger-search-by-creator": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 3},
					Permissions: charon.Permissions{
						charon.UserCanRetrieveAsOwner,
					},
				}, nil)
			},
			req: charonrpc.ListUsersRequest{CreatedBy: qtypes.EqualInt64(2)},
			err: grpcerr.E(codes.PermissionDenied),
		},
	}

	h := listUsersHandler{
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
			defer recoverTest(t)

			actorProviderMock.ExpectedCalls = nil
			userProviderMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.List(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, actorProviderMock, userProviderMock)
		})
	}
}
