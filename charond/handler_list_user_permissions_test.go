package charond

import (
	"context"
	"testing"

	"net"

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
	"google.golang.org/grpc/peer"
)

func TestListUserPermissionsHandler_ListPermissions_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	max := len(charon.AllPermissions)

	_, err := suite.charon.user.SetPermissions(ctx, &charonrpc.SetUserPermissionsRequest{
		UserId:      1,
		Permissions: charon.AllPermissions.Strings(),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("simple", func(t *testing.T) {
		res, err := suite.charon.user.ListPermissions(ctx, &charonrpc.ListUserPermissionsRequest{Id: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Permissions) != max {
			t.Errorf("wrong number of entities, expected %d but got %d", max, len(res.Permissions))
		}
	})
	t.Run("unauthenticated", func(t *testing.T) {
		ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: &net.IPAddr{}})
		_, err := suite.charon.user.ListPermissions(ctx, &charonrpc.ListUserPermissionsRequest{Id: 1})
		assertErrorCode(t, err, codes.Unauthenticated, "mnemosyned: missing access token in metadata")
	})
}

func TestListUserPermissionsHandler_ListPermissions_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	permissionProviderMock := &modelmock.PermissionProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.ListUserPermissionsRequest
		err  error
	}{
		"missing-user-id": {
			init: func(t *testing.T) {
			},
			req: charonrpc.ListUserPermissionsRequest{},
			err: grpcerr.E(codes.InvalidArgument),
		},
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.ListUserPermissionsRequest{Id: 1},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"cannot-list-as-a-stranger-if-missing-permission": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{
							charon.GroupCanRetrieve,
							charon.PermissionCanModify,
							charon.PermissionCanDelete,
							charon.PermissionCanCreate,
						},
						User: &model.UserEntity{ID: 1, IsStaff: true},
					}, nil).
					Once()
			},
			req: charonrpc.ListUserPermissionsRequest{Id: 10},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-list-as-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 10, IsSuperuser: true},
					}, nil).
					Once()
				permissionProviderMock.On("FindByUserID", mock.Anything, int64(1)).
					Return([]*model.PermissionEntity{{
						ID:        1,
						Subsystem: "sub",
						Module:    "mod",
						Action:    "act",
					}}, nil).
					Once()
			},
			req: charonrpc.ListUserPermissionsRequest{Id: 1},
		},
		"can-your-own-groups": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 12},
					}, nil).
					Once()
				permissionProviderMock.On("FindByUserID", mock.Anything, int64(12)).
					Return([]*model.PermissionEntity{{
						ID:        1,
						Subsystem: "sub",
						Module:    "mod",
						Action:    "act",
					}}, nil).
					Once()
			},
			req: charonrpc.ListUserPermissionsRequest{Id: 12},
		},
		"can-list-with-permissions": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserPermissionCanRetrieve},
						User:        &model.UserEntity{ID: 10},
					}, nil).
					Once()
				permissionProviderMock.On("FindByUserID", mock.Anything, int64(1)).
					Return([]*model.PermissionEntity{{
						ID:        1,
						Subsystem: "sub",
						Module:    "mod",
						Action:    "act",
					}}, nil).
					Once()
			},
			req: charonrpc.ListUserPermissionsRequest{Id: 1},
		},
		"storage-query-cancellation": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 1, IsSuperuser: true},
					}, nil).
					Once()
				permissionProviderMock.On("FindByUserID", mock.Anything, int64(10)).
					Return(nil, context.Canceled).
					Once()
			},
			req: charonrpc.ListUserPermissionsRequest{Id: 10},
			err: grpcerr.E(codes.Canceled),
		},
	}

	h := listUserPermissionsHandler{
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
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			actorProviderMock.ExpectedCalls = nil
			permissionProviderMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.ListPermissions(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, actorProviderMock, permissionProviderMock)
		})
	}
}
