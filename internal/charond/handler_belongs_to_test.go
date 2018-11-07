package charond

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
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

func TestBelongsToHandler_BelongsTo_Unit(t *testing.T) {
	userGroupsProviderMock := &modelmock.UserGroupsProvider{}
	actorProviderMock := &sessionmock.ActorProvider{}

	h := belongsToHandler{
		handler: &handler{
			ActorProvider: actorProviderMock,
			repository: repositories{
				userGroups: userGroupsProviderMock,
			},
		},
	}

	cases := map[string]struct {
		init func(t *testing.T)
		req  charonrpc.BelongsToRequest
		err  error
	}{
		"invalid-user-id": {
			init: func(t *testing.T) {},
			req: charonrpc.BelongsToRequest{
				UserId:  -1,
				GroupId: 10,
			},
			err: grpcerr.E(codes.InvalidArgument),
		},
		"invalid-group-id": {
			init: func(t *testing.T) {},
			req: charonrpc.BelongsToRequest{
				UserId:  10,
				GroupId: -1,
			},
			err: grpcerr.E(codes.InvalidArgument),
		},
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.BelongsToRequest{
				UserId:  2,
				GroupId: 5,
			},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"can-check-as-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{
						ID:          1,
						IsSuperuser: true,
					}}, nil).
					Once()
				userGroupsProviderMock.On("Exists", mock.Anything, int64(2), int64(5)).Return(true, nil).
					Once()
			},
			req: charonrpc.BelongsToRequest{
				UserId:  2,
				GroupId: 5,
			},
		},
		"can-check-for-yourself": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{
						ID: 1,
					}}, nil).
					Once()
				userGroupsProviderMock.On("Exists", mock.Anything, int64(1), int64(5)).Return(true, nil).
					Once()
			},
			req: charonrpc.BelongsToRequest{
				UserId:  1,
				GroupId: 5,
			},
		},
		"can-check-with-permission-as-a-stranger": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{
							charon.UserGroupCanCheckBelongingAsStranger,
						},
						User: &model.UserEntity{
							ID: 1,
						},
					}, nil).
					Once()
				userGroupsProviderMock.On("Exists", mock.Anything, int64(2), int64(5)).Return(true, nil).
					Once()
			},
			req: charonrpc.BelongsToRequest{
				UserId:  2,
				GroupId: 5,
			},
		},
		"cannot-check-without-permission": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID: 1,
						},
					}, nil).
					Once()
			},
			req: charonrpc.BelongsToRequest{
				UserId:  2,
				GroupId: 5,
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-check-if-exists-query-fails": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID: 1,
						},
					}, nil).
					Once()
				userGroupsProviderMock.On("Exists", mock.Anything, int64(1), int64(5)).Return(false, context.DeadlineExceeded).
					Once()
			},
			req: charonrpc.BelongsToRequest{
				UserId: 1, GroupId: 5,
			},
			err: grpcerr.E(codes.DeadlineExceeded),
		},
	}

	for hint, c := range cases {

		t.Run(hint, func(t *testing.T) {
			// reset mocks between cases
			actorProviderMock.ExpectedCalls = []*mock.Call{}
			userGroupsProviderMock.ExpectedCalls = []*mock.Call{}

			c.init(t)

			_, err := h.BelongsTo(context.Background(), &c.req)
			if c.err != nil {
				if !grpcerr.Match(c.err, err) {
					t.Fatalf("errors do not match, got '%v'", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}

			if !mock.AssertExpectationsForObjects(t, actorProviderMock, userGroupsProviderMock) {
				t.Errorf("mock expectetions failure")
			}
		})
	}
}

func TestBelongsToHandler_BelongsTo_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	resAct, err := suite.charon.auth.Actor(timeout(ctx), &wrappers.StringValue{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	_, groups := suite.createGroups(t, timeout(ctx))

	_, err = suite.charon.user.SetGroups(ctx, &charonrpc.SetUserGroupsRequest{
		UserId: resAct.Id,
		Groups: groups[:len(groups)/2],
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	cases := map[string]func(t *testing.T){
		"belongs": func(t *testing.T) {
			res, err := suite.charon.auth.BelongsTo(timeout(ctx), &charonrpc.BelongsToRequest{
				UserId:  resAct.Id,
				GroupId: groups[0],
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if !res.Value {
				t.Error("expected to belong")
			}
		},
		"not-belongs": func(t *testing.T) {
			res, err := suite.charon.auth.BelongsTo(timeout(ctx), &charonrpc.BelongsToRequest{
				UserId:  resAct.Id,
				GroupId: groups[len(groups)-1],
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.Value {
				t.Error("expected to not belong")
			}
		},
		"group-does-not-exists": func(t *testing.T) {
			res, err := suite.charon.auth.BelongsTo(timeout(ctx), &charonrpc.BelongsToRequest{
				UserId:  resAct.Id,
				GroupId: 99999999,
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.Value {
				t.Error("expected to not belong")
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}
