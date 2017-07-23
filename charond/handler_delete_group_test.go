package charond

import (
	"context"
	"fmt"
	"math"
	"net"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynetest"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
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
			if grpc.Code(err) != codes.NotFound {
				t.Errorf("wrong status code, expected %s but got %s", codes.NotFound.String(), grpc.Code(err).String())
			}
		},
		"assigned": func(t *testing.T) {
			_, err := suite.charon.group.Delete(timeout(ctx), &charonrpc.DeleteGroupRequest{
				Id: groups[0],
			})
			if grpc.Code(err) != codes.FailedPrecondition {
				t.Errorf("wrong status code, expected %s but got %s", codes.FailedPrecondition.String(), grpc.Code(err).String())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

func TestDeleteGroupHandler_Delete_Unit(t *testing.T) {
	gpm := &model.MockGroupProvider{}
	upm := &model.MockUserProvider{}
	ppm := &model.MockPermissionProvider{}
	sm := &mnemosynetest.SessionManagerClient{}

	sessionOnContext := func(t *testing.T, id int64) {
		sm.On("Context", mock.Anything, &empty.Empty{}, mock.Anything).Return(&mnemosynerpc.ContextResponse{
			Session: &mnemosynerpc.Session{
				SubjectId: session.ActorIDFromInt64(id).String(),
				AccessToken: func() string {
					tkn, err := mnemosyne.RandomAccessToken()
					if err != nil {
						t.Fatalf("token generation error: %s", err.Error())
					}
					return tkn
				}(),
			},
		}, nil).Once()
	}

	h := deleteGroupHandler{
		handler: &handler{
			session: sm,
			repository: repositories{
				user:       upm,
				group:      gpm,
				permission: ppm,
			},
		},
	}

	cases := map[string]func(t *testing.T){
		"invalid-id": func(t *testing.T) {
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: -1})
			assertErrorCode(t, err, codes.InvalidArgument, "group cannot be deleted, invalid id: -1")
		},
		"cannot-remove-if-group-does-not-exists": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.GroupCanDelete.Subsystem(),
					Module:    charon.GroupCanDelete.Module(),
					Action:    charon.GroupCanDelete.Action(),
				},
			}, nil).Once()
			gpm.On("DeleteOneByID", mock.Anything, int64(5)).Return(int64(0), nil).
				Once()
			sessionOnContext(t, 1)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			assertErrorCode(t, err, codes.NotFound, "group cannot be removed, does not exists")
		},
		"cannot-remove-if-group-have-users-assigned": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.GroupCanDelete.Subsystem(),
					Module:    charon.GroupCanDelete.Module(),
					Action:    charon.GroupCanDelete.Action(),
				},
			}, nil).Once()
			gpm.On("DeleteOneByID", mock.Anything, int64(5)).Return(int64(0), &pq.Error{
				Constraint: model.TableUserGroupsConstraintGroupIDForeignKey,
			}).Once()
			sessionOnContext(t, 1)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			assertErrorCode(t, err, codes.FailedPrecondition, "group cannot be removed, users are assigned to it")
		},
		"cannot-remove-if-group-have-permissions-assigned": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.GroupCanDelete.Subsystem(),
					Module:    charon.GroupCanDelete.Module(),
					Action:    charon.GroupCanDelete.Action(),
				},
			}, nil).Once()
			gpm.On("DeleteOneByID", mock.Anything, int64(5)).Return(int64(0), &pq.Error{
				Constraint: model.TableGroupPermissionsConstraintGroupIDForeignKey,
			}).Once()
			sessionOnContext(t, 1)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			assertErrorCode(t, err, codes.FailedPrecondition, "group cannot be removed, permissions are assigned to it")
		},
		"cannot-remove-from-localhost": func(t *testing.T) {
			ctx := peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.TCPAddr{
					IP: net.IPv4(127, 0, 0, 1),
				},
			})
			sm.On("Context", mock.Anything, &empty.Empty{}, mock.Anything).Return(nil, errf(codes.Unknown, "random mnemosyne error")).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			sessionOnContext(t, 1)
			_, err := h.Delete(ctx, &charonrpc.DeleteGroupRequest{Id: 5})
			assertErrorCode(t, err, codes.PermissionDenied, "group cannot be removed from localhost")
		},
		"cannot-remove-if-session-does-not-exists": func(t *testing.T) {
			sm.On("Context", mock.Anything, &empty.Empty{}, mock.Anything).Return(nil, errf(codes.NotFound, "session does not exists")).
				Once()
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			sessionOnContext(t, 1)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			assertErrorCode(t, err, codes.Unauthenticated, "session does not exists")
		},
		"can-remove-if-query-deadline-exceeded": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID:          1,
				IsSuperuser: true,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			gpm.On("DeleteOneByID", mock.Anything, int64(5)).Return(int64(0), context.DeadlineExceeded).
				Once()
			sessionOnContext(t, 1)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			if err != context.DeadlineExceeded {
				t.Fatalf("wrong error, expected %v but got %v", context.DeadlineExceeded, err)
			}
		},
		"cannot-remove-if-missing-permissions": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			sessionOnContext(t, 1)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			assertErrorCode(t, err, codes.PermissionDenied, "group cannot be removed, missing permission")
		},
		"can-remove-as-superuser": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID:          1,
				IsSuperuser: true,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			gpm.On("DeleteOneByID", mock.Anything, int64(5)).Return(int64(1), nil).
				Once()
			sessionOnContext(t, 1)
			_, err := h.Delete(context.Background(), &charonrpc.DeleteGroupRequest{Id: 5})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		},
	}

	for hint, c := range cases {
		// reset mocks between cases
		sm.ExpectedCalls = []*mock.Call{}
		ppm.ExpectedCalls = []*mock.Call{}
		gpm.ExpectedCalls = []*mock.Call{}
		upm.ExpectedCalls = []*mock.Call{}

		t.Run(hint, c)
	}
}
