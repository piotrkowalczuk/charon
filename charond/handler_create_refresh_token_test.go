package charond

import (
	"context"
	"testing"

	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
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

func TestCreateRefreshTokenHandler_Create_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cases := map[string]func(t *testing.T){
		"only-notes": func(t *testing.T) {
			req := &charonrpc.CreateRefreshTokenRequest{
				Notes: &ntypes.String{
					Valid: true,
					Chars: "note",
				},
			}
			res, err := suite.charon.refreshToken.Create(timeout(ctx), req)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if res.RefreshToken.Notes.GetChars() != req.Notes.GetChars() {
				t.Errorf("wrong notes, expected %s but got %s", req.GetNotes(), res.GetRefreshToken().GetNotes())
			}
			if res.RefreshToken.ExpireAt != req.ExpireAt {
				t.Errorf("wrong expire at, expected %#v but got %#v", req.GetExpireAt(), res.GetRefreshToken().GetExpireAt())
			}
		},
	}

	for hint, fn := range cases {
		t.Run(hint, fn)
	}
}

func TestCreateRefreshTokenHandler_Create_Unit(t *testing.T) {
	actorProviderMock := &sessionmock.ActorProvider{}
	refreshTokenProviderMock := &modelmock.RefreshTokenProvider{}

	cases := map[string]struct {
		req  charonrpc.CreateRefreshTokenRequest
		init func(*testing.T)
		err  error
	}{
		"session-does-not-exists": {
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(nil, grpcerr.E(codes.Unauthenticated, "session does not exists")).
					Once()
			},
			req: charonrpc.CreateRefreshTokenRequest{},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"missing-permission": {
			req: charonrpc.CreateRefreshTokenRequest{},
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsStaff: true}}, nil).
					Once()
			},
			err: grpcerr.E(codes.PermissionDenied),
		},
		"broken-date-in-request": {
			req: charonrpc.CreateRefreshTokenRequest{ExpireAt: &timestamp.Timestamp{
				Seconds: -62135596801,
			}},
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
			},
			err: grpcerr.E(codes.InvalidArgument),
		},
		"broken-date-in-response": {
			req: charonrpc.CreateRefreshTokenRequest{},
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				refreshTokenProviderMock.On("Create", mock.Anything, mock.Anything).Return(&model.RefreshTokenEntity{
					ExpireAt: pq.NullTime{Time: brokenDate(), Valid: true},
				}, nil).Once()
			},
			err: grpcerr.E(codes.Internal),
		},
		"creator-does-not-exists": {
			req: charonrpc.CreateRefreshTokenRequest{},
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				refreshTokenProviderMock.On("Create", mock.Anything, mock.Anything).
					Return(nil, &pq.Error{Constraint: model.TableRefreshTokenConstraintCreatedByForeignKey}).
					Once()
			},
			err: grpcerr.E(codes.NotFound),
		},
		"conflict": {
			req: charonrpc.CreateRefreshTokenRequest{},
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				refreshTokenProviderMock.On("Create", mock.Anything, mock.Anything).
					Return(nil, &pq.Error{Constraint: model.TableRefreshTokenConstraintTokenUserIDUnique}).
					Once()
			},
			err: grpcerr.E(codes.AlreadyExists),
		},
		"user-does-not-exists": {
			req: charonrpc.CreateRefreshTokenRequest{},
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				refreshTokenProviderMock.On("Create", mock.Anything, mock.Anything).
					Return(nil, &pq.Error{Constraint: model.TableRefreshTokenConstraintUserIDForeignKey}).
					Once()
			},
			err: grpcerr.E(codes.NotFound),
		},
		"insert-cancel": {
			req: charonrpc.CreateRefreshTokenRequest{},
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				refreshTokenProviderMock.On("Create", mock.Anything, mock.Anything).
					Return(nil, context.Canceled).
					Once()
			},
			err: grpcerr.E(codes.Canceled),
		},
		"can-create-as-superuser": {
			req: charonrpc.CreateRefreshTokenRequest{},
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{User: &model.UserEntity{IsSuperuser: true}}, nil).
					Once()
				refreshTokenProviderMock.On("Create", mock.Anything, mock.Anything).Return(&model.RefreshTokenEntity{
					ExpireAt: pq.NullTime{Time: time.Now(), Valid: true},
				}, nil).Once()
			},
		},
		"can-create-with-permissions": {
			req: charonrpc.CreateRefreshTokenRequest{},
			init: func(t *testing.T) {
				actorProviderMock.On("Actor", mock.Anything).
					Return(&session.Actor{
						Permissions: charon.Permissions{charon.RefreshTokenCanCreate},
						User:        &model.UserEntity{},
					}, nil).
					Once()
				refreshTokenProviderMock.On("Create", mock.Anything, mock.Anything).Return(&model.RefreshTokenEntity{
					ExpireAt: pq.NullTime{Time: time.Now(), Valid: true},
				}, nil).Once()
			},
		},
	}

	h := createRefreshTokenHandler{
		handler: &handler{
			ActorProvider: actorProviderMock,
			repository: repositories{
				refreshToken: refreshTokenProviderMock,
			},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			defer recoverTest(t)

			actorProviderMock.ExpectedCalls = nil
			refreshTokenProviderMock.ExpectedCalls = nil

			c.init(t)

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
			if !mock.AssertExpectationsForObjects(t, actorProviderMock, refreshTokenProviderMock) {
				return
			}
		})
	}
}
