package charond

import (
	"context"
	"errors"
	"testing"

	"database/sql"

	"time"

	"github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/charon/internal/password/passwordmock"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/charon/internal/session/sessionmock"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUserHandler_Create_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	userProviderMock := &modelmock.UserProvider{}
	hasherMock := &passwordmock.Hasher{}

	h := createUserHandler{
		hasher: hasherMock,
		handler: &handler{
			ActorProvider: actorProviderMock,
			repository: repositories{
				user: userProviderMock,
			},
		},
	}

	successInit := func(act session.Actor) func(*testing.T, *charonrpc.CreateUserRequest) {
		return func(t *testing.T, r *charonrpc.CreateUserRequest) {
			actorProviderMock.On("Actor", mock.Anything).
				Return(&act, nil).
				Once()
			hasherMock.On("Hash", []byte(r.PlainPassword)).
				Return([]byte{1, 2, 3}, nil).
				Once()
			userProviderMock.On("Create", mock.Anything, mock.Anything, mock.Anything).
				Return(&model.UserEntity{}, nil).
				Once()
		}
	}

	cases := map[string]struct {
		req  charonrpc.CreateUserRequest
		init func(*testing.T, *charonrpc.CreateUserRequest)
		err  error
	}{
		"missing-username": {
			req: charonrpc.CreateUserRequest{
				Username:      "12",
				PlainPassword: "password",
			},
			init: func(_ *testing.T, _ *charonrpc.CreateUserRequest) {},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		"missing-password": {
			req: charonrpc.CreateUserRequest{
				Username: "username",
			},
			init: func(_ *testing.T, _ *charonrpc.CreateUserRequest) {},
			err:  grpcerr.E(codes.InvalidArgument),
		},
		"missing-permission": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{}}, nil).
					Once()
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"password-hashing-failure": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				hasherMock.On("Hash", []byte(r.PlainPassword)).
					Return(nil, errors.New("example error")).
					Once()
			},
			err: grpcerr.E(codes.Internal),
		},
		"actor-cannot-be-retrieved": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Internal, "example error")).
					Once()
			},
			err: grpcerr.E(codes.Internal),
		},
		"actor-cannot-be-retrieved-and-superuser-already-exists": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsSuperuser:   ntypes.True(),
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Internal, "example error")).
					Once()
				userProviderMock.On("Count", mock.Anything).
					Return(int64(1), nil).
					Once()
			},
			err: grpcerr.E(codes.AlreadyExists),
		},
		"actor-cannot-be-retrieved-and-number-of-users-cannot-be-checked": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsSuperuser:   ntypes.True(),
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Internal, "example error")).
					Once()
				userProviderMock.On("Count", mock.Anything).
					Return(int64(0), sql.ErrConnDone).
					Once()
			},
			err: grpcerr.E(codes.Internal),
		},
		"secure-password-as-regular-user": {
			req: charonrpc.CreateUserRequest{
				Username:       "username",
				SecurePassword: []byte("password"),
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.UserCanCreate},
						User:        &model.UserEntity{},
					}, nil).
					Once()
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"only-superuser-can-create-another-superuser": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsSuperuser:   ntypes.True(),
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{}}, nil).
					Once()
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"having-staff-member-creation-permission-is-not-enough-to-create-superuser": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsSuperuser:   ntypes.True(),
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{
							charon.UserCanCreateStaff,
						},
						User: &model.UserEntity{},
					}, nil).
					Once()
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"staff-member-requires-permission": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsStaff:       ntypes.True(),
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{},
						Permissions: charon.Permissions{
							charon.UserCanCreate,
						},
					}, nil).
					Once()
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"being-staff-member-is-not-enough-to-create-another-staff-member": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsStaff:       ntypes.True(),
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{IsStaff: true},
					}, nil).
					Once()
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"success-as-superuser": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
			},
			init: successInit(session.Actor{User: &model.UserEntity{IsSuperuser: true}}),
		},
		"success-superuser-as-superuser": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsSuperuser:   ntypes.True(),
			},
			init: successInit(session.Actor{User: &model.UserEntity{IsSuperuser: true}}),
		},
		"success-superuser-as-local": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsSuperuser:   ntypes.True(),
			},
			init: successInit(session.Actor{IsLocal: true, User: &model.UserEntity{}}),
		},
		"success-as-user": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
			},
			init: successInit(session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.UserCanCreate,
				},
			}),
		},
		"success-staff-as-user": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsStaff:       ntypes.True(),
			},
			init: successInit(session.Actor{
				User: &model.UserEntity{ID: 2},
				Permissions: charon.Permissions{
					charon.UserCanCreateStaff,
				},
			}),
		},
		"storage-returns-broken-entity": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsStaff:       ntypes.True(),
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 2},
						Permissions: charon.Permissions{
							charon.UserCanCreateStaff,
						},
					}, nil).
					Once()
				hasherMock.On("Hash", []byte(r.PlainPassword)).
					Return([]byte{1, 2, 3}, nil).
					Once()
				userProviderMock.On("Create", mock.Anything, mock.Anything, mock.Anything).
					Return(&model.UserEntity{
						CreatedAt: time.Date(1, 1, 0, 0, 0, 0, 0, time.UTC),
					}, nil).
					Once()
			},
			err: grpcerr.E(codes.Internal),
		},
		"user-with-such-username-already-exists": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsStaff:       ntypes.True(),
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 2},
						Permissions: charon.Permissions{
							charon.UserCanCreateStaff,
						},
					}, nil).
					Once()
				hasherMock.On("Hash", []byte(r.PlainPassword)).
					Return([]byte{1, 2, 3}, nil).
					Once()
				userProviderMock.On("Create", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, &pq.Error{Constraint: model.TableUserConstraintUsernameUnique}).
					Once()
			},
			err: grpcerr.E(codes.AlreadyExists),
		},
		"storage-returns-random-error-on-failure": {
			req: charonrpc.CreateUserRequest{
				Username:      "username",
				PlainPassword: "password",
				IsStaff:       ntypes.True(),
			},
			init: func(t *testing.T, r *charonrpc.CreateUserRequest) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						User: &model.UserEntity{ID: 2},
						Permissions: charon.Permissions{
							charon.UserCanCreateStaff,
						},
					}, nil).
					Once()
				hasherMock.On("Hash", []byte(r.PlainPassword)).
					Return([]byte{1, 2, 3}, nil).
					Once()
				userProviderMock.On("Create", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("example error")).
					Once()
			},
			err: grpcerr.E(codes.Internal),
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			defer recoverTest(t)

			//sessionMock.ExpectedCalls = nil
			//permissionProviderMock.ExpectedCalls = nil
			hasherMock.ExpectedCalls = nil
			userProviderMock.ExpectedCalls = nil
			actorProviderMock.ExpectedCalls = nil

			c.init(t, &c.req)

			_, err := h.Create(context.TODO(), &c.req)
			if c.err != nil {
				if !grpcerr.Match(c.err, err) {
					t.Fatalf("error do not match, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !mock.AssertExpectationsForObjects(t, hasherMock, userProviderMock, actorProviderMock) {
				return
			}
		})
	}
}

func TestCreateUserHandler_Create_E2E(t *testing.T) {
	suite := &endToEndSuite{
		userAgent: "charonctl",
	}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cases := map[string]func(t *testing.T){
		"full": func(t *testing.T) {
			req := &charonrpc.CreateUserRequest{
				Username:      "username-full",
				FirstName:     "first-name-full",
				LastName:      "last-name-full",
				PlainPassword: "plain-password-full",
			}
			res, err := suite.charon.user.Create(timeout(ctx), req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.User.Username != req.Username {
				t.Errorf("wrong username, expected %s but got %s", req.Username, res.User.Username)
			}
			if res.User.FirstName != req.FirstName {
				t.Errorf("wrong first name, expected %#v but got %#v", req.FirstName, res.User.FirstName)
			}
			if res.User.LastName != req.LastName {
				t.Errorf("wrong last name, expected %#v but got %#v", req.LastName, res.User.LastName)
			}
		},
		"superuser-twice": func(t *testing.T) {
			_, err := suite.charon.user.Create(timeout(ctx), &charonrpc.CreateUserRequest{
				Username:      "superuser-1",
				FirstName:     "first-name-1",
				LastName:      "last-name-1",
				PlainPassword: "plain-password-1",
				IsSuperuser:   ntypes.True(),
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			_, err = suite.charon.user.Create(timeout(ctx), &charonrpc.CreateUserRequest{
				Username:      "superuser-2",
				FirstName:     "first-name-2",
				LastName:      "last-name-2",
				PlainPassword: "plain-password-2",
				IsSuperuser:   ntypes.True(),
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		},
		"local": func(t *testing.T) {
			req := &charonrpc.CreateUserRequest{
				Username:      "username-local",
				FirstName:     "first-name-local",
				LastName:      "last-name-local",
				PlainPassword: "plain-password-local",
			}
			_, err := suite.charon.user.Create(timeout(context.Background()), req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		},
		"same-username-twice": func(t *testing.T) {
			req := &charonrpc.CreateUserRequest{
				Username:      "username-same-username-twice",
				FirstName:     "first-name-same-username-twice",
				LastName:      "last-name-same-username-twice",
				PlainPassword: "plain-password-same-username-twice",
			}
			_, err := suite.charon.user.Create(timeout(ctx), req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			_, err = suite.charon.user.Create(timeout(ctx), req)
			if status.Code(err) != codes.AlreadyExists {
				t.Errorf("wrong status code, expected %s but got %s", codes.AlreadyExists.String(), status.Code(err).String())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}
