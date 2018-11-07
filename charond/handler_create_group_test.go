package charond

import (
	"context"
	"testing"
	"time"

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
	"google.golang.org/grpc/status"
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
			if status.Code(err) != codes.AlreadyExists {
				t.Fatalf("wrong status code, expected %s but got %s", codes.AlreadyExists.String(), status.Code(err).String())
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
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("wrong status code, expected %s but got %s", codes.InvalidArgument.String(), status.Code(err).String())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

func TestCreateGroupHandler_Create_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	groupProviderMock := &modelmock.GroupProvider{}

	h := createGroupHandler{
		handler: &handler{
			ActorProvider: actorProviderMock,
			repository: repositories{
				group: groupProviderMock,
			},
		},
	}

	cases := map[string]struct {
		init func(t *testing.T)
		req  charonrpc.CreateGroupRequest
		err  error
	}{
		"name-to-short": {
			init: func(t *testing.T) {},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		"can-create-with-permission": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.GroupCanCreate},
						User: &model.UserEntity{
							ID: 1,
						},
					}, nil).
					Once()
				groupProviderMock.On("Create", mock.Anything, int64(1), "name", mock.AnythingOfType("*ntypes.String")).Return(&model.GroupEntity{
					CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
					Name:      "name",
				}, nil).Once()
			},
			req: charonrpc.CreateGroupRequest{
				Name: "name",
			},
		},
		"cannot-persist-if-such-name-already-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID:          1,
							IsSuperuser: true,
						},
					}, nil).
					Once()
				groupProviderMock.On("Create", mock.Anything, int64(1), "name", mock.AnythingOfType("*ntypes.String")).
					Return(nil, &pq.Error{Constraint: model.TableGroupConstraintNameUnique}).Once()
			},
			req: charonrpc.CreateGroupRequest{
				Name: "name",
			},
			err: grpcerr.E(codes.AlreadyExists),
		},
		"can-create-as-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID:          1,
							IsSuperuser: true,
						},
					}, nil).
					Once()
				groupProviderMock.On("Create", mock.Anything, int64(1), "name", mock.AnythingOfType("*ntypes.String")).
					Return(&model.GroupEntity{
						CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
						Name:      "name",
					}, nil).Once()
			},
			req: charonrpc.CreateGroupRequest{
				Name: "name",
			},
		},
		"cannot-reply-if-entity-is-broken": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID:          1,
							IsSuperuser: true,
						},
					}, nil).
					Once()
				groupProviderMock.On("Create", mock.Anything, int64(1), "name", mock.AnythingOfType("*ntypes.String")).
					Return(&model.GroupEntity{
						CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
						Name:      "name",
						CreatedAt: time.Date(1, 1, 0, 0, 0, 0, 0, time.UTC),
					}, nil).Once()
			},
			req: charonrpc.CreateGroupRequest{
				Name: "name",
			},
			err: grpcerr.E(codes.Internal),
		},
		"cannot-create-without-permission": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID:      1,
							IsStaff: true,
						},
					}, nil).
					Once()
			},
			req: charonrpc.CreateGroupRequest{
				Name: "name",
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-check-if-exists-query-fails": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID:          1,
							IsSuperuser: true,
						},
					}, nil).
					Once()
				groupProviderMock.On("Create", mock.Anything, int64(1), "name", mock.AnythingOfType("*ntypes.String")).Return(nil, context.DeadlineExceeded).
					Once()
			},
			req: charonrpc.CreateGroupRequest{
				Name: "name",
			},
			err: grpcerr.E(codes.DeadlineExceeded),
		},
		"cannot-check-if-session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.NotFound, "session not found")).
					Once()
			},
			req: charonrpc.CreateGroupRequest{
				Name: "name",
			},
			err: grpcerr.E(codes.NotFound),
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			// reset mocks between cases
			actorProviderMock.ExpectedCalls = []*mock.Call{}
			groupProviderMock.ExpectedCalls = []*mock.Call{}

			c.init(t)

			_, err := h.Create(context.Background(), &c.req)
			if c.err != nil {
				if !grpcerr.Match(c.err, err) {
					t.Fatalf("errors do not match, got '%v'", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}

			mock.AssertExpectationsForObjects(t, actorProviderMock, groupProviderMock)
		})

	}
}
