package charond

import (
	"context"
	"database/sql"
	"math"
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
)

func TestDeleteUserHandler_Delete_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)
	_, err := suite.charon.auth.Actor(timeout(ctx), &wrappers.StringValue{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	_, users := suite.createUsers(t, timeout(ctx))
	_, groups := suite.createGroups(t, timeout(ctx))

	_, err = suite.charon.user.SetGroups(timeout(ctx), &charonrpc.SetUserGroupsRequest{
		UserId: users[1],
		Groups: groups[:len(groups)/2],
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	cases := map[string]func(t *testing.T){
		"not-assigned": func(t *testing.T) {
			done, err := suite.charon.user.Delete(timeout(ctx), &charonrpc.DeleteUserRequest{
				Id: users[len(users)-1],
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if !done.Value {
				t.Error("group expected to be removed")
			}
		},
		"not-existing": func(t *testing.T) {
			_, err := suite.charon.user.Delete(timeout(ctx), &charonrpc.DeleteUserRequest{
				Id: math.MaxInt64,
			})
			assertErrorCode(t, err, codes.NotFound, "user does not exists")
		},
		"groups-assigned": func(t *testing.T) {
			_, err := suite.charon.user.Delete(timeout(ctx), &charonrpc.DeleteUserRequest{
				Id: users[1],
			})
			assertErrorCode(t, err, codes.FailedPrecondition, "user cannot be removed, groups are assigned to it")
		},
		"permissions-assigned": func(t *testing.T) {
			t.Skip("TODO: implement")

			_, err := suite.charon.user.Delete(timeout(ctx), &charonrpc.DeleteUserRequest{
				Id: users[1],
			})
			assertErrorCode(t, err, codes.FailedPrecondition, "user cannot be removed, permissions are assigned to it")
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

func TestDeleteUserHandler_Delete_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	userProviderMock := &modelmock.UserProvider{}

	h := deleteUserHandler{
		handler: &handler{
			ActorProvider: actorProviderMock,
			repository: repositories{
				user: userProviderMock,
			},
		},
	}

	cases := map[string]struct {
		init func(*testing.T, *charonrpc.DeleteUserRequest)
		req  charonrpc.DeleteUserRequest
		err  error
	}{
		"invalid-id": {
			init: func(_ *testing.T, _ *charonrpc.DeleteUserRequest) {},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		"cannot-remove-if-actor-does-not-exists": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.PermissionDenied, "actor not found")).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 10},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-remove-from-localhost": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						IsLocal: true,
						User: &model.UserEntity{
							ID: 1,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{ID: 11}, nil).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-remove-delete-query-timeout": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteAsStranger},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{ID: 11}, nil).
					Once()
				userProviderMock.On("DeleteOneByID", mock.Anything, int64(11)).Return(int64(0), context.DeadlineExceeded).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.DeadlineExceeded),
		},
		"cannot-remove-find-user-query-timeout": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteAsStranger},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(nil, context.DeadlineExceeded).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.DeadlineExceeded),
		},
		"cannot-remove-itself": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID: 10,
						},
						Permissions: charon.Permissions{
							charon.UserCanDeleteAsStranger,
							charon.UserCanDeleteAsOwner,
							charon.UserCanDeleteStaffAsStranger,
							charon.UserCanDeleteStaffAsOwner,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 10},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-remove-itself-even-if-superuser": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID:          10,
							IsSuperuser: true,
						},
						Permissions: charon.Permissions{
							charon.UserCanDeleteAsStranger,
							charon.UserCanDeleteAsOwner,
							charon.UserCanDeleteStaffAsStranger,
							charon.UserCanDeleteStaffAsOwner,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{
					ID: 10, IsSuperuser: true,
				}, nil).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 10},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-remove-superuser": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID: 10,
						},
						Permissions: charon.Permissions{
							charon.UserCanDeleteAsStranger,
							charon.UserCanDeleteAsOwner,
							charon.UserCanDeleteStaffAsStranger,
							charon.UserCanDeleteStaffAsOwner,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10, IsSuperuser: true}, nil).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 10},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-remove-as-superuser": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{
							ID:          11,
							IsSuperuser: true,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(10)).Return(&model.UserEntity{ID: 10}, nil).
					Once()
				userProviderMock.On("DeleteOneByID", mock.Anything, int64(10)).Return(int64(1), nil).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 10},
		},
		"cannot-remove-as-stranger": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteAsOwner},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
					ID:        11,
					CreatedBy: ntypes.Int64{Int64: 12, Valid: true},
				}, nil).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-remove-as-owner": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteAsStranger},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
					CreatedBy: ntypes.Int64{Int64: 10, Valid: true},
					ID:        11,
				}, nil).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"cannot-remove-if-user-have-permissions-assigned": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteAsOwner},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
					CreatedBy: ntypes.Int64{Int64: 10, Valid: true},
					ID:        11,
				}, nil).
					Once()
				userProviderMock.On("DeleteOneByID", mock.Anything, int64(11)).Return(int64(0), &pq.Error{
					Constraint: model.TableUserPermissionsConstraintUserIDForeignKey,
				}).Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.FailedPrecondition),
		},
		"cannot-remove-if-user-have-groups-assigned": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteAsOwner},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
					CreatedBy: ntypes.Int64{Int64: 10, Valid: true},
					ID:        11,
				}, nil).
					Once()
				userProviderMock.On("DeleteOneByID", mock.Anything, int64(11)).Return(int64(0), &pq.Error{
					Constraint: model.TableUserGroupsConstraintUserIDForeignKey,
				}).Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.FailedPrecondition),
		},
		"can-delete-as-stranger-but-does-not-exists": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteAsStranger},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(nil, sql.ErrNoRows).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.NotFound),
		},
		"can-delete-as-stranger-but-not-a-superuser": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteAsStranger},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
					ID:          11,
					IsSuperuser: true,
				}, nil)
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-delete-as-stranger-but-not-a-staff-member": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteAsStranger},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
					ID:      11,
					IsStaff: true,
				}, nil)
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-delete-staff-member-but-not-as-a-stranger": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteStaffAsOwner},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
					ID:      11,
					IsStaff: true,
				}, nil)
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-delete-staff-member-but-not-as-a-owner": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteStaffAsStranger},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).Return(&model.UserEntity{
					ID:        11,
					IsStaff:   true,
					CreatedBy: ntypes.Int64{Int64: 10, Valid: true},
				}, nil)
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"can-delete-as-owner": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteAsOwner},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).
					Return(&model.UserEntity{
						ID:        11,
						CreatedBy: ntypes.Int64{Int64: 10, Valid: true},
					}, nil)
				userProviderMock.On("DeleteOneByID", mock.Anything, int64(11)).
					Return(int64(1), nil).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
		},
		"can-delete-staff-member-as-owner": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteStaffAsOwner},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).
					Return(&model.UserEntity{
						ID:        11,
						IsStaff:   true,
						CreatedBy: ntypes.Int64{Int64: 10, Valid: true},
					}, nil)
				userProviderMock.On("DeleteOneByID", mock.Anything, int64(11)).
					Return(int64(1), nil).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
		},
		"can-delete-staff-member-as-stranger": {
			init: func(t *testing.T, r *charonrpc.DeleteUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanDeleteStaffAsStranger},
						User: &model.UserEntity{
							ID: 10,
						},
					}, nil).
					Once()
				userProviderMock.On("FindOneByID", mock.Anything, int64(11)).
					Return(&model.UserEntity{
						ID:      11,
						IsStaff: true,
					}, nil)
				userProviderMock.On("DeleteOneByID", mock.Anything, int64(11)).
					Return(int64(1), nil).
					Once()
			},
			req: charonrpc.DeleteUserRequest{Id: 11},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			// reset mocks between cases
			actorProviderMock.ExpectedCalls = []*mock.Call{}
			userProviderMock.ExpectedCalls = []*mock.Call{}

			c.init(t, &c.req)

			_, err := h.Delete(context.Background(), &c.req)
			if c.err != nil {
				if !grpcerr.Match(c.err, err) {
					t.Fatalf("errors do not match, got '%v'", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}

			mock.AssertExpectationsForObjects(t, actorProviderMock, userProviderMock)
		})
	}
}
