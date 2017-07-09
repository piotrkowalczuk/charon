package charond

import (
	"context"
	"testing"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogoutHandler_Logout(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	tkn, err := suite.charon.auth.Login(context.TODO(), &charonrpc.LoginRequest{
		Username: "test",
		Password: "test",
	})
	if err != nil {
		t.Fatalf("unexpected login error: %s: with code %s", grpc.ErrorDesc(err), grpc.Code(err))
	}

	if _, err := suite.charon.auth.Logout(context.TODO(), &charonrpc.LogoutRequest{
		AccessToken: tkn.Value,
	}); err != nil {
		t.Errorf("logout failure: %s", err.Error())
	}
	ok, err := suite.charon.auth.IsAuthenticated(context.TODO(), &charonrpc.IsAuthenticatedRequest{
		AccessToken: tkn.Value,
	})
	if err != nil {
		t.Fatalf("unexpected is authenticated error: %s: with code %s", grpc.ErrorDesc(err), grpc.Code(err))
	}
	if ok.Value {
		t.Errorf("user should not be authenticated")
	}
}

func TestLogoutHandler_Logout_missingToken(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	_, err := suite.charon.auth.Logout(context.TODO(), &charonrpc.LogoutRequest{})
	if err == nil {
		t.Fatal("error should not be nil")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.InvalidArgument {
			t.Errorf("wrong status code, expected %s but got %s", codes.InvalidArgument.String(), st.Code().String())
		}
	} else {
		t.Error("wrong error type")
	}
}
