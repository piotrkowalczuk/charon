package charond

import (
	"context"
	"fmt"
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

func TestListGroupsHandler_List_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	max := 5
	for i := 0; i < max; i++ {
		_, err := suite.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{
			Name: fmt.Sprintf("group-%d", i),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("simple", func(t *testing.T) {
		res, err := suite.charon.group.List(ctx, &charonrpc.ListGroupsRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Groups) != max {
			t.Errorf("wrong number of entities, expected %d but got %d", max, len(res.Groups))
		}
	})
	t.Run("unauthenticated", func(t *testing.T) {
		ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: &net.IPAddr{}})
		_, err := suite.charon.group.List(ctx, &charonrpc.ListGroupsRequest{})
		assertErrorCode(t, err, codes.Unauthenticated, "mnemosyned: missing access token in metadata")
	})
	t.Run("offset", func(t *testing.T) {
		res, err := suite.charon.group.List(ctx, &charonrpc.ListGroupsRequest{
			Offset: ntypes.NewInt64(int64(max - 1)),
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Groups) != 1 {
			t.Errorf("wrong number of entities, expected %d but got %d", 1, len(res.Groups))
		}
	})
	t.Run("order-by", func(t *testing.T) {
		res, err := suite.charon.group.List(ctx, &charonrpc.ListGroupsRequest{
			OrderBy: []*charonrpc.Order{
				{
					Name:       "name",
					Descending: true,
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		exp := fmt.Sprintf("group-%d", max-1)
		if res.Groups[0].Name != exp {
			t.Errorf("wrong group name, expected %s but got %s", exp, res.Groups[0].Name)
		}
	})
}

func TestListGroupsHandler_List_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	groupProviderMock := &modelmock.GroupProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.ListGroupsRequest
		err  error
	}{
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.ListGroupsRequest{},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"cannot-list-as-a-stranger-if-missing-permission": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{
							charon.PermissionCanRetrieve,
							charon.GroupCanModify,
							charon.GroupCanDelete,
							charon.GroupCanCreate,
						},
						User: &model.UserEntity{ID: 1, IsStaff: true},
					}, nil).
					Once()
			},
			req: charonrpc.ListGroupsRequest{},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-list-as-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 1, IsSuperuser: true},
					}, nil).
					Once()
				groupProviderMock.On("Find", mock.Anything, &model.GroupFindExpr{
					Limit:   10,
					Offset:  0,
					OrderBy: []model.RowOrder{},
				}).Return([]*model.GroupEntity{{
					ID:   1,
					Name: "example",
				}}, nil)
			},
			req: charonrpc.ListGroupsRequest{},
		},
		"can-list-with-permissions": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.GroupCanRetrieve},
						User:        &model.UserEntity{ID: 1},
					}, nil).
					Once()
				groupProviderMock.On("Find", mock.Anything, &model.GroupFindExpr{
					Limit:   10,
					Offset:  0,
					OrderBy: []model.RowOrder{},
				}).Return([]*model.GroupEntity{{
					ID:   1,
					Name: "example",
				}}, nil)
			},
			req: charonrpc.ListGroupsRequest{},
		},
		"reverse-mapping-failure": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.GroupCanRetrieve},
						User:        &model.UserEntity{ID: 1},
					}, nil).
					Once()
				groupProviderMock.On("Find", mock.Anything, &model.GroupFindExpr{
					Limit:   10,
					Offset:  0,
					OrderBy: []model.RowOrder{},
				}).Return([]*model.GroupEntity{{
					ID:        1,
					Name:      "example",
					CreatedAt: brokenDate(),
				}}, nil)
			},
			req: charonrpc.ListGroupsRequest{},
			err: grpcerr.E(codes.Internal),
		},
		"storage-query-cancellation": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.GroupCanRetrieve},
						User:        &model.UserEntity{ID: 1},
					}, nil).
					Once()
				groupProviderMock.On("Find", mock.Anything, &model.GroupFindExpr{
					Limit:   10,
					Offset:  0,
					OrderBy: []model.RowOrder{},
				}).Return(nil, context.Canceled)
			},
			req: charonrpc.ListGroupsRequest{},
			err: grpcerr.E(codes.Canceled),
		},
	}

	h := listGroupsHandler{
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

			_, err := h.List(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, actorProviderMock, groupProviderMock)
		})
	}
}
