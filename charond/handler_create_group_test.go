package charond

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
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
)

func TestCreateGroupHandler_Create_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cases := map[string]func(t *testing.T){
		"full": func(t *testing.T) {
			req := &charonrpc.CreateGroupRequest{
				Name: "name-full",
				Description: &ntypes.String{
					Valid: true,
					Chars: "description",
				},
			}
			res, err := suite.charon.group.Create(timeout(ctx), req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.Group.Name != req.Name {
				t.Errorf("wrong name, expected %s but got %s", req.Name, res.Group.Name)
			}
			if res.Group.Description != req.Description.StringOr("") {
				t.Errorf("wrong description, expected %#v but got %#v", req.Description.StringOr(""), res.Group.Description)
			}
		},
		"only-name": func(t *testing.T) {
			req := &charonrpc.CreateGroupRequest{
				Name: "name-only-name",
			}
			res, err := suite.charon.group.Create(timeout(ctx), req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.Group.Name != req.Name {
				t.Errorf("wrong name, expected %s but got %s", req.Name, res.Group.Name)
			}
			if res.Group.Description != "" {
				t.Errorf("wrong description, expected %#v but got %#v", "", res.Group.Description)
			}
		},
		"same-name-twice": func(t *testing.T) {
			req := &charonrpc.CreateGroupRequest{
				Name: "same-name-twice",
			}
			_, err := suite.charon.group.Create(timeout(ctx), req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			_, err = suite.charon.group.Create(timeout(ctx), req)
			if grpc.Code(err) != codes.AlreadyExists {
				t.Fatalf("wrong status code, expected %s but got %s", codes.AlreadyExists.String(), grpc.Code(err).String())
			}
		},
		"only-description": func(t *testing.T) {
			req := &charonrpc.CreateGroupRequest{
				Description: &ntypes.String{
					Valid: true,
					Chars: "description",
				},
			}
			_, err := suite.charon.group.Create(timeout(ctx), req)
			if grpc.Code(err) != codes.InvalidArgument {
				t.Fatalf("wrong status code, expected %s but got %s", codes.InvalidArgument.String(), grpc.Code(err).String())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

func TestCreateGroupHandler_Create_Unit(t *testing.T) {
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

	h := createGroupHandler{
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
		"name-to-short": func(t *testing.T) {
			_, err := h.Create(context.Background(), &charonrpc.CreateGroupRequest{
				Name: "12",
			})
			assertErrorCode(t, err, codes.InvalidArgument, "group name is required and needs to be at least 3 characters long")
		},
		"can-create-with-permission": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{
				{
					Subsystem: charon.GroupCanCreate.Subsystem(),
					Module:    charon.GroupCanCreate.Module(),
					Action:    charon.GroupCanCreate.Action(),
				},
			}, nil).Once()
			gpm.On("Create", mock.Anything, int64(1), "name", mock.AnythingOfType("*ntypes.String")).Return(&model.GroupEntity{
				CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
				Name:      "name",
			}, nil).Once()
			sessionOnContext(t, 1)
			_, err := h.Create(context.Background(), &charonrpc.CreateGroupRequest{
				Name: "name",
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		},
		"can-create-as-superuser": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID:          1,
				IsSuperuser: true,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			gpm.On("Create", mock.Anything, int64(1), "name", mock.AnythingOfType("*ntypes.String")).
				Return(&model.GroupEntity{
					CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
					Name:      "name",
				}, nil).Once()
			sessionOnContext(t, 1)
			_, err := h.Create(context.Background(), &charonrpc.CreateGroupRequest{
				Name: "name",
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		},
		"cannot-reply-if-entity-is-broken": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID:          1,
				IsSuperuser: true,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			gpm.On("Create", mock.Anything, int64(1), "name", mock.AnythingOfType("*ntypes.String")).
				Return(&model.GroupEntity{
					CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
					Name:      "name",
					CreatedAt: time.Date(1, 1, 0, 0, 0, 0, 0, time.UTC),
				}, nil).Once()
			sessionOnContext(t, 1)
			_, err := h.Create(context.Background(), &charonrpc.CreateGroupRequest{
				Name: "name",
			})
			assertErrorCode(t, err, codes.Internal, "group entity mapping failure: timestamp: seconds:-62135683200  before 0001-01-01")
		},
		"cannot-create-without-permission": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID: 1,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			sessionOnContext(t, 1)
			_, err := h.Create(context.Background(), &charonrpc.CreateGroupRequest{
				Name: "name",
			})
			assertErrorCode(t, err, codes.PermissionDenied, "group cannot be created, missing permission")
		},
		"cannot-check-if-exists-query-fails": func(t *testing.T) {
			upm.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
				ID:          1,
				IsSuperuser: true,
			}, nil)
			ppm.On("FindByUserID", mock.Anything, int64(1)).Return([]*model.PermissionEntity{}, nil).
				Once()
			gpm.On("Create", mock.Anything, int64(1), "name", mock.AnythingOfType("*ntypes.String")).Return(nil, context.DeadlineExceeded).
				Once()
			sessionOnContext(t, 1)
			_, err := h.Create(context.Background(), &charonrpc.CreateGroupRequest{
				Name: "name",
			})
			if err != context.DeadlineExceeded {
				t.Fatalf("wrong error, expected %v but got %v", context.DeadlineExceeded, err.Error())
			}
		},
		"cannot-check-if-session-does-not-exists": func(t *testing.T) {
			sm.On("Context", mock.Anything, &empty.Empty{}, mock.Anything).Return(nil, context.Canceled).Once()
			_, err := h.Create(context.Background(), &charonrpc.CreateGroupRequest{
				Name: "name",
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
		gpm.ExpectedCalls = []*mock.Call{}
		upm.ExpectedCalls = []*mock.Call{}

		t.Run(hint, c)
	}
}

func TestCreateGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charonrpc.CreateGroupRequest
		act session.Actor
	}{
		{
			req: charonrpc.CreateGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.GroupCanCreate,
				},
			},
		},
		{
			req: charonrpc.CreateGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &createGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestCreateGroupHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.CreateGroupRequest
		act session.Actor
	}{
		{
			req: charonrpc.CreateGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 2},
			},
		},
		{
			req: charonrpc.CreateGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charonrpc.CreateGroupRequest{},
			act: session.Actor{
				User: &model.UserEntity{ID: 1},
			},
		},
	}

	h := &createGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
