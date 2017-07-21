package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func testRegisterPermissionsHandlerRegister(t *testing.T, permissions []string, fn func(t *testing.T, err error)) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_, err := suite.charon.permission.Register(ctx, &charonrpc.RegisterPermissionsRequest{
		Permissions: permissions,
	})
	fn(t, err)
}

func TestRegisterPermissionsHandler_Register(t *testing.T) {
	testRegisterPermissionsHandlerRegister(t, []string{
		"a:b:c",
		"a:bb:cc",
	}, func(t *testing.T, err error) {
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	})
}

func TestRegisterPermissionsHandler_Register_inconsistentSubsystem(t *testing.T) {
	testRegisterPermissionsHandlerRegister(t, []string{
		"a:b:c",
		"aa:bb:cc",
	}, func(t *testing.T, err error) {
		if err == nil {
			t.Fatal("missing error")
		}
		if st, ok := status.FromError(err); ok {
			if st.Code() != codes.InvalidArgument {
				t.Errorf("wrong error code, expected %s but got %s", codes.InvalidArgument.String(), st.Code().String())
			}
		}
	})
}

func TestRegisterPermissionsHandler_Register_emptySubsystem(t *testing.T) {
	testRegisterPermissionsHandlerRegister(t, []string{
		":b:c",
	}, func(t *testing.T, err error) {
		if err == nil {
			t.Fatal("missing error")
		}
		if st, ok := status.FromError(err); ok {
			if st.Code() != codes.InvalidArgument {
				t.Errorf("wrong error code, expected %s but got %s", codes.InvalidArgument.String(), st.Code().String())
			}
		}
	})
}
