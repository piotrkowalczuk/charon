package charond

import (
	"context"
	"database/sql"
	"testing"

	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model/modelmock"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynetest"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestIsAuthenticatedHandler_IsAuthenticated_E2E(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("metadata is missing in context")
	}
	tkn, ok := md[mnemosyne.AccessTokenMetadataKey]
	if !ok {
		t.Fatalf("access token is missing in metadata")
	}
	res, err := suite.charon.auth.IsAuthenticated(ctx, &charonrpc.IsAuthenticatedRequest{AccessToken: tkn[0]})
	if err != nil {
		t.Fatal(err)
	}
	if !res.GetValue() {
		t.Error("expected true")
	}
	res, err = suite.charon.auth.IsAuthenticated(ctx, &charonrpc.IsAuthenticatedRequest{AccessToken: "131231231312321"})
	if err != nil {
		t.Fatal(err)
	}
	if res.GetValue() {
		t.Error("expected false")
	}
}

func TestIsAuthenticatedHandler_IsAuthenticated_Unit(t *testing.T) {
	sessionMock := &mnemosynetest.SessionManagerClient{}
	userProviderMock := &modelmock.UserProvider{}

	res := &mnemosynerpc.GetResponse{
		Session: &mnemosynerpc.Session{
			AccessToken:   "123",
			SubjectId:     "charon:user:1",
			SubjectClient: "test",
		},
	}
	cases := map[string]struct {
		req  charonrpc.IsAuthenticatedRequest
		init func(*testing.T, *charonrpc.IsAuthenticatedRequest)
		err  error
	}{
		"success": {
			req: charonrpc.IsAuthenticatedRequest{AccessToken: "access-token"},
			init: func(t *testing.T, r *charonrpc.IsAuthenticatedRequest) {
				sessionMock.On("Get", mock.Anything, &mnemosynerpc.GetRequest{
					AccessToken: r.AccessToken,
				}, mock.Anything).Return(res, nil)
				userProviderMock.On("Exists", mock.Anything, int64(1)).Return(true, nil)
			},
		},
		"session-not-found": {
			req: charonrpc.IsAuthenticatedRequest{AccessToken: "access-token"},
			init: func(t *testing.T, r *charonrpc.IsAuthenticatedRequest) {
				sessionMock.On("Get", mock.Anything, &mnemosynerpc.GetRequest{
					AccessToken: r.AccessToken,
				}, mock.Anything).Return(nil, status.Errorf(codes.NotFound, "session not found"))
			},
		},
		"session-canceled": {
			req: charonrpc.IsAuthenticatedRequest{AccessToken: "access-token"},
			init: func(t *testing.T, r *charonrpc.IsAuthenticatedRequest) {
				sessionMock.On("Get", mock.Anything, &mnemosynerpc.GetRequest{
					AccessToken: r.AccessToken,
				}, mock.Anything).Return(nil, status.Errorf(codes.Canceled, "something went wrong"))
			},
			err: grpcerr.E(codes.Canceled),
		},
		"session-store-returns-broken-subject-id": {
			req: charonrpc.IsAuthenticatedRequest{AccessToken: "access-token"},
			init: func(t *testing.T, r *charonrpc.IsAuthenticatedRequest) {
				sessionMock.On("Get", mock.Anything, &mnemosynerpc.GetRequest{
					AccessToken: r.AccessToken,
				}, mock.Anything).Return(&mnemosynerpc.GetResponse{
					Session: &mnemosynerpc.Session{
						AccessToken:   "123",
						SubjectId:     "1",
						SubjectClient: "test",
					},
				}, nil)
			},
			err: grpcerr.E(codes.Internal),
		},
		"user-does-not-exists": {
			req: charonrpc.IsAuthenticatedRequest{AccessToken: "access-token"},
			init: func(t *testing.T, r *charonrpc.IsAuthenticatedRequest) {
				sessionMock.On("Get", mock.Anything, &mnemosynerpc.GetRequest{
					AccessToken: r.AccessToken,
				}, mock.Anything).Return(res, nil)
				userProviderMock.On("Exists", mock.Anything, int64(1)).Return(false, sql.ErrNoRows)
			},
		},
		"user-existance-check-timeout": {
			req: charonrpc.IsAuthenticatedRequest{AccessToken: "access-token"},
			init: func(t *testing.T, r *charonrpc.IsAuthenticatedRequest) {
				sessionMock.On("Get", mock.Anything, &mnemosynerpc.GetRequest{
					AccessToken: r.AccessToken,
				}, mock.Anything).Return(res, nil)
				userProviderMock.On("Exists", mock.Anything, int64(1)).Return(false, context.DeadlineExceeded)
			},
			err: grpcerr.E(codes.DeadlineExceeded),
		},
		"no-access-token": {
			req:  charonrpc.IsAuthenticatedRequest{},
			init: func(t *testing.T, r *charonrpc.IsAuthenticatedRequest) {},
			err:  grpcerr.E(codes.InvalidArgument, "authentication status cannot be checked, missing access token"),
		},
	}

	h := &isAuthenticatedHandler{
		handler: &handler{
			session: sessionMock,
			repository: repositories{
				user: userProviderMock,
			},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			sessionMock.ExpectedCalls = nil
			userProviderMock.ExpectedCalls = nil

			c.init(t, &c.req)

			_, err := h.IsAuthenticated(context.TODO(), &c.req)
			assertError(t, c.err, err)

			mock.AssertExpectationsForObjects(t, sessionMock, userProviderMock)
		})
	}
}
