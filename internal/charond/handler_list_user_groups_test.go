package charond

import (
	"context"
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
)

func TestListUserGroupsHandler_ListGroups_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)
	userID := int64(1)

	createGroupResp, err := suite.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{
		Name: "existing-group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = suite.charon.user.SetGroups(ctx, &charonrpc.SetUserGroupsRequest{
		UserId: userID,
		Groups: []int64{createGroupResp.GetGroup().GetId()},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	listGroupsResp, err := suite.charon.user.ListGroups(ctx, &charonrpc.ListUserGroupsRequest{
		Id: userID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if len(listGroupsResp.GetGroups()) != 1 {
		t.Errorf("wrong number of groups, expected 1 got %d", len(listGroupsResp.GetGroups()))
	}
}

func TestListUserGroupsHandler_ListGroups_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	groupProviderMock := &modelmock.GroupProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.ListUserGroupsRequest
		err  error
	}{
		"missing-user-id": {
			init: func(t *testing.T) {
			},
			req: charonrpc.ListUserGroupsRequest{},
			err: grpcerr.E(codes.InvalidArgument),
		},
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.ListUserGroupsRequest{Id: 1},
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
			req: charonrpc.ListUserGroupsRequest{Id: 10},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-list-as-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 10, IsSuperuser: true},
					}, nil).
					Once()
				groupProviderMock.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.GroupEntity{{
					ID:   1,
					Name: "example",
				}}, nil)
			},
			req: charonrpc.ListUserGroupsRequest{Id: 1},
		},
		"can-your-own-groups": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 12},
					}, nil).
					Once()
				groupProviderMock.On("FindByUserID", mock.Anything, int64(12)).Return([]*model.GroupEntity{{
					ID:   1,
					Name: "example",
				}}, nil)
			},
			req: charonrpc.ListUserGroupsRequest{Id: 12},
		},
		"can-list-with-permissions": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserGroupCanRetrieve},
						User:        &model.UserEntity{ID: 10},
					}, nil).
					Once()
				groupProviderMock.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.GroupEntity{{
					ID:   1,
					Name: "example",
				}}, nil)
			},
			req: charonrpc.ListUserGroupsRequest{Id: 1},
		},
		"reverse-mapping-failure": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserGroupCanRetrieve},
						User:        &model.UserEntity{ID: 10},
					}, nil).
					Once()
				groupProviderMock.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.GroupEntity{{
					ID:        1,
					Name:      "example",
					CreatedAt: brokenDate(),
				}}, nil)
			},
			req: charonrpc.ListUserGroupsRequest{Id: 1},
			err: grpcerr.E(codes.Internal),
		},
		"storage-query-cancellation": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 1, IsSuperuser: true},
					}, nil).
					Once()
				groupProviderMock.On("FindByUserID", mock.Anything, int64(10)).Return(nil, context.Canceled)
			},
			req: charonrpc.ListUserGroupsRequest{Id: 10},
			err: grpcerr.E(codes.Canceled),
		},
	}

	h := listUserGroupsHandler{
		handler: &handler{
			logger:        sklog.NewTestLogger(t),
			ActorProvider: actorProviderMock,
			repository: repositories{
				group: groupProviderMock,
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
			groupProviderMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.ListGroups(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, actorProviderMock, groupProviderMock)
		})
	}
}
