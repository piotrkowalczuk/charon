package charond

import (
	"context"
	"fmt"
	"math"
	"net"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/charon/internal/session/sessionmock"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func TestDeleteGroupHandler_Delete_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)
	resAct, err := suite.charon.auth.Actor(timeout(ctx), &wrappers.StringValue{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	var groups []int64
	for i := 0; i < 10; i++ {
		resGroup, err := suite.charon.group.Create(timeout(ctx), &charonrpc.CreateGroupRequest{
			Name: fmt.Sprintf("name-%d", i),
			Description: &ntypes.String{
				Valid: true,
				Chars: fmt.Sprintf("description-%d", i),
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		groups = append(groups, resGroup.Group.Id)
	}

	_, err = suite.charon.user.SetGroups(timeout(ctx), &charonrpc.SetUserGroupsRequest{
		UserId: resAct.Id,
		Groups: groups[:len(groups)/2],
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	cases := map[string]func(t *testing.T){
		"not-assigned": func(t *testing.T) {
			done, err := suite.charon.group.Delete(timeout(ctx), &charonrpc.DeleteGroupRequest{
				Id: groups[len(groups)-1],
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if !done.Value {
				t.Error("group expected to be removed")
			}
		},
		"not-existing": func(t *testing.T) {
			_, err := suite.charon.group.Delete(timeout(ctx), &charonrpc.DeleteGroupRequest{
				Id: math.MaxInt64,
			})
			if status.Code(err) != codes.NotFound {
				t.Errorf("wrong status code, expected %s but got %s", codes.NotFound.String(), status.Code(err).String())
			}
		},
		"assigned": func(t *testing.T) {
			_, err := suite.charon.group.Delete(timeout(ctx), &charonrpc.DeleteGroupRequest{
				Id: groups[0],
			})
			if status.Code(err) != codes.FailedPrecondition {
				t.Errorf("wrong status code, expected %s but got %s", codes.FailedPrecondition.String(), status.Code(err).String())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

func TestDeleteGroupHandler_Delete_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	groupProviderMock := &modelmock.GroupProvider{}

	h := deleteGroupHandler{
		handler: &handler{
			ActorProvider: actorProviderMock,
			repository: repositories{
				group: groupProviderMock,
			},
		},
	}

	cases := map[string]func(t *testing.T){
		"invalid-id": func(t *testing.T) {
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: -1})
			if !grpcerr.Match(grpcerr.E(codes.InvalidArgument), err) {
				t.Fatalf("errors do not match, got '%v'", err)
			}
		},
		"cannot-remove-if-group-does-not-exists": func(t *testing.T) {
			actorProviderMock.On("Actor", mock.Anything).
				Return(&session.Actor{
					Permissions: charon.Permissions{charon.GroupCanDelete},
					User: &model.UserEntity{
						ID: 1,
					},
				}, nil).
				Once()
			groupProviderMock.On("DeleteOneByID", mock.Anything, int64(5)).
				Return(int64(0), nil).
				Once()
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			if !grpcerr.Match(grpcerr.E(codes.NotFound), err) {
				t.Fatalf("errors do not match, got '%v'", err)
			}
		},
		"cannot-remove-if-group-have-users-assigned": func(t *testing.T) {
			actorProviderMock.On("Actor", mock.Anything).
				Return(&session.Actor{
					Permissions: charon.Permissions{charon.GroupCanDelete},
					User: &model.UserEntity{
						ID: 1,
					},
				}, nil).
				Once()
			groupProviderMock.On("DeleteOneByID", mock.Anything, int64(5)).
				Return(int64(0), &pq.Error{
					Constraint: model.TableUserGroupsConstraintGroupIDForeignKey,
				}).
				Once()
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			if !grpcerr.Match(grpcerr.E(codes.FailedPrecondition), err) {
				t.Fatalf("errors do not match, got '%v'", err)
			}
		},
		"cannot-remove-if-group-have-permissions-assigned": func(t *testing.T) {
			actorProviderMock.On("Actor", mock.Anything).
				Return(&session.Actor{
					Permissions: charon.Permissions{charon.GroupCanDelete},
					User: &model.UserEntity{
						ID: 1,
					},
				}, nil).
				Once()
			groupProviderMock.On("DeleteOneByID", mock.Anything, int64(5)).
				Return(int64(0), &pq.Error{
					Constraint: model.TableGroupPermissionsConstraintGroupIDForeignKey,
				}).
				Once()
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			if !grpcerr.Match(grpcerr.E(codes.FailedPrecondition), err) {
				t.Fatalf("errors do not match, got '%v'", err)
			}
		},
		"cannot-remove-from-localhost": func(t *testing.T) {
			ctx := peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.TCPAddr{
					IP: net.IPv4(127, 0, 0, 1),
				},
			})
			actorProviderMock.On("Actor", mock.Anything).
				Return(&session.Actor{
					IsLocal: true,
					User: &model.UserEntity{
						ID: 1,
					},
				}, nil).
				Once()
			_, err := h.Delete(ctx, &charonrpc.DeleteGroupRequest{Id: 5})
			if !grpcerr.Match(grpcerr.E(codes.PermissionDenied), err) {
				t.Fatalf("errors do not match, got '%v'", err)
			}
		},
		"cannot-remove-if-session-does-not-exists": func(t *testing.T) {
			actorProviderMock.On("Actor", mock.Anything).
				Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
				Once()
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			if !grpcerr.Match(grpcerr.E(codes.Unauthenticated), err) {
				t.Fatalf("errors do not match, got '%v'", err)
			}
		},
		"can-remove-if-query-deadline-exceeded": func(t *testing.T) {
			actorProviderMock.On("Actor", mock.Anything).
				Return(&session.Actor{
					User: &model.UserEntity{
						ID:          1,
						IsSuperuser: true,
					},
				}, nil).
				Once()
			groupProviderMock.On("DeleteOneByID", mock.Anything, int64(5)).Return(int64(0), context.DeadlineExceeded).
				Once()
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			if !grpcerr.Match(grpcerr.E(codes.DeadlineExceeded), err) {
				t.Fatalf("errors do not match, got '%v'", err)
			}
		},
		"cannot-remove-if-missing-permissions": func(t *testing.T) {
			actorProviderMock.On("Actor", mock.Anything).
				Return(&session.Actor{
					User: &model.UserEntity{
						ID: 1,
					},
				}, nil).
				Once()
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			if !grpcerr.Match(grpcerr.E(codes.PermissionDenied), err) {
				t.Fatalf("errors do not match, got '%v'", err)
			}
		},
		"can-remove-as-superuser": func(t *testing.T) {
			actorProviderMock.On("Actor", mock.Anything).
				Return(&session.Actor{
					User: &model.UserEntity{
						ID:          1,
						IsSuperuser: true,
					},
				}, nil).
				Once()
			groupProviderMock.On("DeleteOneByID", mock.Anything, int64(5)).Return(int64(1), nil).
				Once()
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		},
	}

	for hint, c := range cases {
		// reset mocks between cases
		actorProviderMock.ExpectedCalls = nil
		groupProviderMock.ExpectedCalls = []*mock.Call{}

		t.Run(hint, c)

		if !mock.AssertExpectationsForObjects(t, actorProviderMock, groupProviderMock) {
			return
		}
	}
}
