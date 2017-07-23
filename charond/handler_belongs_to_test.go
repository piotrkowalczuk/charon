package charond

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynetest"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
)

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

func TestBelongsToHandler_BelongsTo_Unit(t *testing.T) {
	ugpm := &model.MockUserGroupsProvider{}
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

	h := belongsToHandler{
		handler: &handler{
			session: sm,
			repository: repositories{
				user:       upm,
				userGroups: ugpm,
				permission: ppm,
			},
		},
	}

	cases := map[string]func(t *testing.T){
		"invalid-user-id": func(t *testing.T) {
			_, err := h.BelongsTo(context.Background(), &charonrpc.BelongsToRequest{
				UserId:  -1,
				GroupId: 10,
			})
			assertErrorCode(t, err, codes.InvalidArgument, "user id needs to be greater than zero")
		},
		"invalid-group-id": func(t *testing.T) {
			_, err := h.BelongsTo(context.Background(), &charonrpc.BelongsToRequest{
				UserId:  10,
				GroupId: -1,
			})
			assertErrorCode(t, err, codes.InvalidArgument, "group id needs to be greater than zero")
		},
		"can-check-as-superuser": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID:          1,
				IsSuperuser: true,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			ugpm.On("Exists", mock.Anything, int64(2), int64(5)).Return(true, nil).
				Once()
			sessionOnContext(t, 1)
			_, err := h.BelongsTo(context.Background(), &charonrpc.BelongsToRequest{
				UserId:  2,
				GroupId: 5,
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		},
		"can-check-for-yourself": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			ugpm.On("Exists", mock.Anything, int64(1), int64(5)).Return(true, nil).
				Once()
			sessionOnContext(t, 1)
			_, err := h.BelongsTo(context.Background(), &charonrpc.BelongsToRequest{
				UserId:  1,
				GroupId: 5,
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		},
		"can-check-with-permission-as-a-stranger": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.UserGroupCanCheckBelongingAsStranger.Subsystem(),
					Module:    charon.UserGroupCanCheckBelongingAsStranger.Module(),
					Action:    charon.UserGroupCanCheckBelongingAsStranger.Action(),
				},
			}, nil).
				Once()
			ugpm.On("Exists", mock.Anything, int64(2), int64(5)).Return(true, nil).
				Once()
			sessionOnContext(t, 1)
			_, err := h.BelongsTo(context.Background(), &charonrpc.BelongsToRequest{
				UserId:  2,
				GroupId: 5,
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		},
		"cannot-check-without-permission": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			ugpm.On("Exists", mock.Anything, int64(2), int64(5)).Return(true, nil).
				Once()
			sessionOnContext(t, 1)
			_, err := h.BelongsTo(context.Background(), &charonrpc.BelongsToRequest{
				UserId:  2,
				GroupId: 5,
			})
			assertErrorCode(t, err, codes.PermissionDenied, "group belonging cannot be checked, missing permission")
		},
		"cannot-check-if-exists-query-fails": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			ugpm.On("Exists", mock.Anything, int64(1), int64(5)).Return(false, context.DeadlineExceeded).
				Once()
			sessionOnContext(t, 1)
			_, err := h.BelongsTo(context.Background(), &charonrpc.BelongsToRequest{
				UserId: 1, GroupId: 5,
			})
			if err != context.DeadlineExceeded {
				t.Fatalf("wrong error, expected %v but got %v", context.DeadlineExceeded, err.Error())
			}
		},
		"cannot-check-if-session-does-not-exists": func(t *testing.T) {
			sm.On("Context", mock.Anything, &empty.Empty{}, mock.Anything).Return(nil, context.Canceled).Once()
			_, err := h.BelongsTo(context.Background(), &charonrpc.BelongsToRequest{
				UserId: 1, GroupId: 5,
			})
			if err != context.Canceled {
				t.Fatalf("wrong error, expected %v but got %v", context.Canceled, err.Error())
			}
		},
	}

	for hint, c := range cases {
		// reset mocks between cases
		sm.ExpectedCalls = []*mock.Call{}
		ppm.ExpectedCalls = []*mock.Call{}
		ugpm.ExpectedCalls = []*mock.Call{}
		upm.ExpectedCalls = []*mock.Call{}

		t.Run(hint, c)
	}
}

func TestBelongsToHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.BelongsToRequest
		act session.Actor
	}{
		{
			req: charonrpc.BelongsToRequest{UserId: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.UserGroupCanCheckBelongingAsStranger,
				},
			},
		},

		{
			req: charonrpc.BelongsToRequest{UserId: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},

		{
			req: charonrpc.BelongsToRequest{UserId: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
	}

	h := &belongsToHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestBelongsToHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.BelongsToRequest
		act session.Actor
	}{
		{
			req: charonrpc.BelongsToRequest{UserId: 1},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
			},
		},
		{
			req: charonrpc.BelongsToRequest{UserId: 1},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &belongsToHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
