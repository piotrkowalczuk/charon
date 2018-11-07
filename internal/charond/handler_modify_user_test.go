package charond

import (
	"context"
	"database/sql"
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
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestModifyUserHandler_Modify_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	res, err := suite.charon.group.Create(ctx, &charonrpc.CreateGroupRequest{Name: "example"})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("missing-id", func(t *testing.T) {
		_, err := suite.charon.user.Modify(ctx, &charonrpc.ModifyUserRequest{})
		assertErrorCode(t, err, codes.InvalidArgument, "user cannot be modified, invalid id")
	})
	t.Run("ok", func(t *testing.T) {
		_, err := suite.charon.user.Modify(ctx, &charonrpc.ModifyUserRequest{
			Id:        res.Group.Id,
			FirstName: ntypes.NewString("A"),
			LastName:  ntypes.NewString("B"),
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestModifyUserHandler_Modify_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	userProviderMock := &modelmock.UserProvider{}

	cases := map[string]struct {
		init func(*testing.T)
		req  charonrpc.ModifyUserRequest
		err  error
	}{
		"user-id-missing": {
			init: func(t *testing.T) {
			},
			req: charonrpc.ModifyUserRequest{},
			err: grpcerr.E(codes.InvalidArgument),
		},
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"storage-query-cancel": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1), mock.Anything).Return(nil, context.Canceled)
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
			err: grpcerr.E(codes.Canceled),
		},
		"reverse-mapping-failure": {
			init: func(t *testing.T) {
				ent := &model.UserEntity{
					ID:        1,
					FirstName: "name",
					CreatedAt: brokenDate(),
				}
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).Return(ent, nil)
				userProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).Return(ent, nil)
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
			err: grpcerr.E(codes.Internal),
		},
		"not-found": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 2, IsSuperuser: true},
					}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(nil, sql.ErrNoRows)
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
			err: grpcerr.E(codes.NotFound),
		},
		"can-modify-as-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{
						ID:        1,
						FirstName: "123",
						CreatedAt: time.Now(),
					}, nil).
					Once()
				userProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).
					Return(&model.UserEntity{
						ID: 1, FirstName: "123", IsStaff: true, CreatedAt: time.Now(),
					}, nil).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
		},
		"such-username-already-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{ID: 1}, nil).
					Once()
				userProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).
					Return(nil, &pq.Error{Constraint: model.TableUserConstraintUsernameUnique}).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, Username: ntypes.NewString("123")},
			err: grpcerr.E(codes.AlreadyExists),
		},
		"user-deleted-before-update": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{ID: 1}, nil).
					Once()
				userProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).
					Return(nil, sql.ErrNoRows).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, Username: ntypes.NewString("123")},
			err: grpcerr.E(codes.NotFound),
		},
		"update-query-canceled": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2, IsSuperuser: true},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{ID: 1}, nil).
					Once()
				userProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).
					Return(nil, context.Canceled).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, Username: ntypes.NewString("123")},
			err: grpcerr.E(codes.Canceled),
		},
		"cannot-modify-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 2},
						Permissions: charon.Permissions{
							charon.UserCanModifyAsOwner,
							charon.UserCanModifyAsStranger,
							charon.UserCanModifyStaffAsOwner,
							charon.UserCanModifyStaffAsStranger,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
					ID:          1,
					IsSuperuser: true,
					FirstName:   "123",
					CreatedAt:   time.Now(),
				}, nil).Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-promote-to-superuser": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 2},
						Permissions: charon.Permissions{
							charon.UserCanModifyAsOwner,
							charon.UserCanModifyAsStranger,
							charon.UserCanModifyStaffAsOwner,
							charon.UserCanModifyStaffAsStranger,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).Return(&model.UserEntity{
					ID: 1,
				}, nil).Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, IsSuperuser: ntypes.True()},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-modify-as-stranger": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User:        &model.UserEntity{ID: 2},
					Permissions: charon.Permissions{charon.UserCanModifyAsStranger},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{
						ID:        1,
						FirstName: "123",
						CreatedAt: time.Now(),
					}, nil).
					Once()
				userProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).
					Return(&model.UserEntity{
						ID: 1, FirstName: "123", IsStaff: true, CreatedAt: time.Now(),
					}, nil).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
		},
		"cannot-modify-as-stranger": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2},
					Permissions: charon.Permissions{
						charon.UserCanModifyAsOwner,
						charon.UserCanModifyStaffAsOwner,
						charon.UserCanModifyStaffAsStranger,
					},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{
						ID:        1,
						FirstName: "123",
						CreatedAt: time.Now(),
					}, nil).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-modify-as-owner": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 1},
						Permissions: charon.Permissions{
							charon.UserCanModifyStaffAsOwner,
							charon.UserCanModifyStaffAsStranger,
						},
					}, nil).Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{
						ID:        1,
						CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
					}, nil).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-modify-as-owner": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User:        &model.UserEntity{ID: 1},
						Permissions: charon.Permissions{charon.UserCanModifyAsOwner},
					}, nil).Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{
						ID:        1,
						FirstName: "123",
						CreatedAt: time.Now(),
						CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
					}, nil).
					Once()
				userProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).
					Return(&model.UserEntity{
						ID: 1, FirstName: "123", IsStaff: true, CreatedAt: time.Now(),
					}, nil).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
		},
		"can-modify-staff-as-stranger": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User:        &model.UserEntity{ID: 2},
						Permissions: charon.Permissions{charon.UserCanModifyStaffAsStranger},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{
						ID:        1,
						FirstName: "123",
						IsStaff:   true,
						CreatedAt: time.Now(),
						CreatedBy: ntypes.Int64{Int64: 10, Valid: true},
					}, nil).
					Once()
				userProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).
					Return(&model.UserEntity{
						ID: 1, FirstName: "123", IsStaff: true, CreatedAt: time.Now(),
					}, nil).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
		},
		"cannot-modify-staff-as-stranger": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User:        &model.UserEntity{ID: 2},
						Permissions: charon.Permissions{charon.UserCanModifyStaffAsOwner},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{
						ID:        1,
						FirstName: "123",
						IsStaff:   true,
						CreatedBy: ntypes.Int64{Int64: 10, Valid: true},
					}, nil).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-modify-staff-as-owner": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User:        &model.UserEntity{ID: 1},
					Permissions: charon.Permissions{charon.UserCanModifyStaffAsOwner},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{
						ID:        1,
						FirstName: "123",
						IsStaff:   true,
						CreatedAt: time.Now(),
						CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
					}, nil).
					Once()
				userProviderMock.On("UpdateOneByID", mock.Anything, int64(1), mock.Anything).Return(&model.UserEntity{
					ID: 1, FirstName: "123", IsStaff: true, CreatedAt: time.Now(),
				}, nil).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
		},
		"cannot-modify-staff-as-owner": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 1},
					Permissions: charon.Permissions{
						charon.UserCanModifyAsOwner,
						charon.UserCanModifyAsStranger,
					},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{
						ID:        1,
						FirstName: "123",
						IsStaff:   true,
						CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
					}, nil).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-promote-staff": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 1},
					Permissions: charon.Permissions{
						charon.UserCanModifyAsOwner,
						charon.UserCanModifyAsStranger,
						charon.UserCanModifyStaffAsOwner,
						charon.UserCanModifyStaffAsStranger,
					},
				}, nil)
				userProviderMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{
						ID:        1,
						CreatedBy: ntypes.Int64{Int64: 1, Valid: true},
					}, nil).
					Once()
			},
			req: charonrpc.ModifyUserRequest{Id: 1, IsStaff: ntypes.True()},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-modify-without-permissions": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).Return(&session.Actor{
					User: &model.UserEntity{ID: 2},
				}, nil)
			},
			req: charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("123")},
			err: grpcerr.E(codes.PermissionDenied),
		},
	}

	h := modifyUserHandler{
		handler: &handler{
			logger:        zap.L(),
			ActorProvider: actorProviderMock,
			repository: repositories{
				user: userProviderMock,
			},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			defer recoverTest(t)

			actorProviderMock.ExpectedCalls = nil
			userProviderMock.ExpectedCalls = nil

			c.init(t)

			_, err := h.Modify(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, actorProviderMock, userProviderMock)
		})
	}
}

func TestModifyUserHandler_Modify(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cres := testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John",
		LastName:      "Snow",
	})
	_, err := suite.charon.user.Modify(ctx, &charonrpc.ModifyUserRequest{
		Id:        cres.User.Id,
		Username:  ntypes.NewString("john88@snow.com"),
		FirstName: ntypes.NewString("john"),
		LastName:  ntypes.NewString("snow"),
		IsActive:  ntypes.True(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestModifyUserHandler_Modify_nonExistingUser(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_ = testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John",
		LastName:      "Snow",
	})
	_, err := suite.charon.user.Modify(ctx, &charonrpc.ModifyUserRequest{
		Id:        1000,
		Username:  ntypes.NewString("john88@snow.com"),
		FirstName: ntypes.NewString("john"),
		LastName:  ntypes.NewString("snow"),
		IsActive:  ntypes.True(),
	})
	if err == nil {
		t.Fatal("missing error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.NotFound {
			t.Errorf("wrong error code, expected %s but got %s", codes.NotFound.String(), st.Code().String())
		}
	}
}

func TestModifyUserHandler_Modify_wrongID(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_ = testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John",
		LastName:      "Snow",
	})
	_, err := suite.charon.user.Modify(ctx, &charonrpc.ModifyUserRequest{
		Id: -1,
	})
	if err == nil {
		t.Fatal("missing error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.InvalidArgument {
			t.Errorf("wrong error code, expected %s but got %s", codes.InvalidArgument.String(), st.Code().String())
		}
	}
}

func TestModifyUserHandler_Modify_usernameAlreadyExists(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_ = testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John",
		LastName:      "Snow",
	})
	cres := testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john2@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John2",
		LastName:      "Snow2",
	})
	_, err := suite.charon.user.Modify(ctx, &charonrpc.ModifyUserRequest{
		Id:       cres.User.Id,
		Username: ntypes.NewString("john@snow.com"),
	})
	if err == nil {
		t.Fatal("missing error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.AlreadyExists {
			t.Errorf("wrong error code, expected %s but got %s", codes.AlreadyExists.String(), st.Code().String())
		}
	}
}
