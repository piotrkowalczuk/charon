package charond

import (
	"context"
	"sort"
	"testing"

	"net"

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
	"google.golang.org/grpc/peer"
)

func TestListPermissionsHandler_List_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	max := len(charon.AllPermissions)

	t.Run("simple", func(t *testing.T) {
		res, err := suite.charon.permission.List(ctx, &charonrpc.ListPermissionsRequest{})
		if err != nil {
			t.Fatal(err)
		}
		m := max
		if m > 10 {
			m = 10
		}
		if len(res.Permissions) != m {
			t.Errorf("wrong number of entities, expected %d but got %d", m, len(res.Permissions))
		}
	})
	t.Run("unauthenticated", func(t *testing.T) {
		ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: &net.IPAddr{}})
		_, err := suite.charon.permission.List(ctx, &charonrpc.ListPermissionsRequest{})
		assertErrorCode(t, err, codes.Unauthenticated, "mnemosyned: missing access token in metadata")
	})
	t.Run("offset", func(t *testing.T) {
		res, err := suite.charon.permission.List(ctx, &charonrpc.ListPermissionsRequest{
			Offset: ntypes.NewInt64(int64(max - 1)),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Permissions) != 1 {
			t.Errorf("wrong number of entities, expected %d but got %d:\n%v", 1, len(res.Permissions), res.Permissions)
		}
	})
	t.Run("order-by", func(t *testing.T) {
		res, err := suite.charon.permission.List(ctx, &charonrpc.ListPermissionsRequest{
			Limit: ntypes.NewInt64(10000),
			OrderBy: []*charonrpc.Order{
				{
					Name:       model.TablePermissionColumnSubsystem,
					Descending: true,
				},
				{
					Name:       model.TablePermissionColumnModule,
					Descending: true,
				},
				{
					Name:       model.TablePermissionColumnAction,
					Descending: true,
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		exp := charon.AllPermissions
		sort.Sort(sort.Reverse(exp))

		if res.Permissions[0] != exp[0].String() {
			t.Errorf("wrong group name, expected %s but got %s", exp[0].String(), res.Permissions[0])
		}
	})
}

func TestListPermissionsHandler_List_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	permissionProviderMock := &modelmock.PermissionProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.ListPermissionsRequest
		err  error
	}{
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.ListPermissionsRequest{},
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
			req: charonrpc.ListPermissionsRequest{},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-list-as-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 1, IsSuperuser: true},
					}, nil).
					Once()
				permissionProviderMock.On("Find", mock.Anything, &model.PermissionFindExpr{
					Where:   &model.PermissionCriteria{},
					Limit:   10,
					Offset:  0,
					OrderBy: []model.RowOrder{},
				}).Return([]*model.PermissionEntity{{
					ID:        1,
					Subsystem: "subsystem",
					Module:    "module",
					Action:    "action",
				}}, nil)
			},
			req: charonrpc.ListPermissionsRequest{},
		},
		"can-list-with-permissions": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.PermissionCanRetrieve},
						User:        &model.UserEntity{ID: 1},
					}, nil).
					Once()
				permissionProviderMock.On("Find", mock.Anything, &model.PermissionFindExpr{
					Where:   &model.PermissionCriteria{},
					Limit:   10,
					Offset:  0,
					OrderBy: []model.RowOrder{},
				}).Return([]*model.PermissionEntity{{
					ID:        1,
					Subsystem: "subsystem",
					Module:    "module",
					Action:    "action",
				}}, nil)
			},
			req: charonrpc.ListPermissionsRequest{},
		},
		"storage-query-cancellation": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.PermissionCanRetrieve},
						User:        &model.UserEntity{ID: 1},
					}, nil).
					Once()
				permissionProviderMock.On("Find", mock.Anything, &model.PermissionFindExpr{
					Where:   &model.PermissionCriteria{},
					Limit:   10,
					Offset:  0,
					OrderBy: []model.RowOrder{},
				}).Return(nil, context.Canceled)
			},
			req: charonrpc.ListPermissionsRequest{},
			err: grpcerr.E(codes.Canceled),
		},
	}

	h := listPermissionsHandler{
		handler: &handler{
			logger:        zap.L(),
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

			_, err := h.List(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, actorProviderMock, permissionProviderMock)
		})
	}
}

func TestListPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.ListPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.ListPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
				Permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
				},
			},
		},
		{
			req: charonrpc.ListPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &listPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestListPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.ListPermissionsRequest
		act session.Actor
	}{
		{
			req: charonrpc.ListPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.ListPermissionsRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &listPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
