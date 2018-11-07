package charond

import (
	"context"
	"testing"

	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynetest"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestActorHandler_Actor_Unit(t *testing.T) {
	sessionMock := &mnemosynetest.SessionManagerClient{}
	userMock := &modelmock.UserProvider{}
	permissionMock := &modelmock.PermissionProvider{}

	cases := map[string]struct {
		fn  func(*testing.T)
		tok string
		err error
	}{
		"session-through-context": {
			fn: func(t *testing.T) {
				sessionMock.On("Context", mock.Anything, mock.Anything).
					Return(&mnemosynerpc.ContextResponse{
						Session: &mnemosynerpc.Session{
							AccessToken:  "123",
							SubjectId:    "charon:user:1",
							RefreshToken: "abc",
						},
					}, nil).
					Once()
				userMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{ID: 1}, nil).
					Once()
				permissionMock.On("FindByUserID", mock.Anything, int64(1)).
					Return([]*model.PermissionEntity{
						{
							Subsystem: charon.GroupCanRetrieve.Subsystem(),
							Module:    charon.GroupCanRetrieve.Module(),
							Action:    charon.GroupCanRetrieve.Action(),
						},
					}, nil).
					Once()
			},
		},
		"session-through-context-missing-subject-id": {
			fn: func(t *testing.T) {
				sessionMock.On("Context", mock.Anything, mock.Anything).
					Return(&mnemosynerpc.ContextResponse{
						Session: &mnemosynerpc.Session{
							AccessToken:  "123",
							RefreshToken: "abc",
						},
					}, nil).Once()
			},
			err: grpcerr.E(codes.Internal),
		},
		"session-through-value": {
			tok: "123",
			fn: func(t *testing.T) {
				sessionMock.On("Get", mock.Anything, mock.Anything).
					Return(&mnemosynerpc.GetResponse{
						Session: &mnemosynerpc.Session{
							AccessToken:  "123",
							SubjectId:    "charon:user:1",
							RefreshToken: "abc",
						},
					}, nil).
					Once()
				userMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{ID: 1}, nil).
					Once()
				permissionMock.On("FindByUserID", mock.Anything, int64(1)).
					Return([]*model.PermissionEntity{}, nil).
					Once()
			},
		},
		"session-through-value-user-does-not-exists": {
			tok: "123",
			fn: func(t *testing.T) {
				sessionMock.On("Get", mock.Anything, mock.Anything).
					Return(&mnemosynerpc.GetResponse{
						Session: &mnemosynerpc.Session{
							AccessToken:  "123",
							SubjectId:    "charon:user:1",
							RefreshToken: "abc",
						},
					}, nil).
					Once()
				userMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(nil, sql.ErrNoRows).
					Once()
			},
			err: grpcerr.E(codes.NotFound),
		},
		"session-through-value-user-context-canceled": {
			tok: "123",
			fn: func(t *testing.T) {
				sessionMock.On("Get", mock.Anything, mock.Anything).
					Return(&mnemosynerpc.GetResponse{
						Session: &mnemosynerpc.Session{
							AccessToken:  "123",
							SubjectId:    "charon:user:1",
							RefreshToken: "abc",
						},
					}, nil).
					Once()
				userMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(nil, context.Canceled).
					Once()
			},
			err: grpcerr.E(codes.Canceled),
		},
		"session-through-value-user-fetch-postgres-error": {
			tok: "123",
			fn: func(t *testing.T) {
				sessionMock.On("Get", mock.Anything, mock.Anything).
					Return(&mnemosynerpc.GetResponse{
						Session: &mnemosynerpc.Session{
							AccessToken:  "123",
							SubjectId:    "charon:user:1",
							RefreshToken: "abc",
						},
					}, nil).
					Once()
				userMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(nil, &pq.Error{Message: "example"}).
					Once()
			},
			err: grpcerr.E(codes.Internal),
		},
		"session-through-value-permission-context-canceled": {
			tok: "123",
			fn: func(t *testing.T) {
				sessionMock.On("Get", mock.Anything, mock.Anything).
					Return(&mnemosynerpc.GetResponse{
						Session: &mnemosynerpc.Session{
							AccessToken:  "123",
							SubjectId:    "charon:user:1",
							RefreshToken: "abc",
						},
					}, nil).
					Once()
				userMock.On("FindOneByID", mock.Anything, int64(1)).
					Return(&model.UserEntity{ID: 1}, nil).
					Once()
				permissionMock.On("FindByUserID", mock.Anything, int64(1)).
					Return(nil, context.Canceled).
					Once()
			},
			err: grpcerr.E(codes.Canceled),
		},
		"session-not-found-by-context": {
			fn: func(t *testing.T) {
				sessionMock.On("Context", mock.Anything, mock.Anything).
					Return(nil, status.Errorf(codes.NotFound, "session not found")).
					Once()
			},
			err: grpcerr.E(codes.Unauthenticated),
		},
		"session-not-found-by-value": {
			tok: "123",
			fn: func(t *testing.T) {
				sessionMock.On("Get", mock.Anything, mock.Anything).
					Return(nil, status.Errorf(codes.NotFound, "session not found")).
					Once()
			},
			err: grpcerr.E(codes.Unauthenticated),
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			sessionMock.ExpectedCalls = nil
			userMock.ExpectedCalls = nil
			permissionMock.ExpectedCalls = nil

			c.fn(t)

			h := actorHandler{
				handler: &handler{
					session: sessionMock,
					repository: repositories{
						user:       userMock,
						permission: permissionMock,
					},
				},
			}

			_, err := h.Actor(context.TODO(), &wrappers.StringValue{Value: c.tok})
			if !mock.AssertExpectationsForObjects(t, sessionMock, userMock, permissionMock) {
				return
			}
			if c.err != nil {
				if !grpcerr.Match(c.err, err) {
					t.Fatalf("errors do not match, expected '%v', but got '%v'", c.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		})
	}
}

func TestActorHandler_Actor_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("metadata not present in context")
	}

	tok := &wrappers.StringValue{Value: md[mnemosyne.AccessTokenMetadataKey][0]}

	res, err := suite.charon.auth.Actor(context.Background(), tok)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if res.Username != "test" {
		t.Errorf("wrong username, expected %s but got %s", "test", res.Username)
	}

	_, err = suite.charon.user.SetPermissions(ctx, &charonrpc.SetUserPermissionsRequest{
		UserId: res.Id,
		Permissions: []string{
			charon.PermissionCanCreate.String(),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	_, err = suite.charon.auth.Actor(context.Background(), tok)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}
