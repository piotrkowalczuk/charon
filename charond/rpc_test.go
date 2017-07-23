package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestRPCServer_minimal(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)
	permissions := []string{
		"winterfell:castle:can enter as a lord",
		"winterfell:castle:can close as a lord",
	}

	createUserResponse := testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John",
		LastName:      "Snow",
	})
	createGroupResponse := testRPCServerCreateGroup(t, suite, timeout(ctx), &charonrpc.CreateGroupRequest{
		Name: "winterfell",
	})
	registerPermissionsResponse := testRPCServerRegisterPermissions(t, suite, timeout(ctx), &charonrpc.RegisterPermissionsRequest{
		Permissions: permissions,
	})
	if registerPermissionsResponse.Created != 2 {
		t.Fatalf("wrong number of registered Permissions, expected 2 but got %d", registerPermissionsResponse.Created)
	}
	_ = testRPCServerSetUserPermissions(t, suite, timeout(ctx), &charonrpc.SetUserPermissionsRequest{
		UserId:      createUserResponse.User.Id,
		Permissions: permissions,
	})
	_ = testRPCServerSetUserGroups(t, suite, timeout(ctx), &charonrpc.SetUserGroupsRequest{
		UserId: createUserResponse.User.Id,
		Groups: []int64{createGroupResponse.Group.Id},
	})
}

func testRPCServerLogin(t *testing.T, suite *endToEndSuite) context.Context {
	token, err := suite.charon.auth.Login(context.TODO(), &charonrpc.LoginRequest{
		Username: "test",
		Password: "test",
	})
	if err != nil {
		t.Fatalf("unexpected login error: %s: with code %s", grpc.ErrorDesc(err), grpc.Code(err))
	}

	return metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(mnemosyne.AccessTokenMetadataKey, token.Value),
	)
}

func testRPCServerCreateUser(t *testing.T, suite *endToEndSuite, ctx context.Context, req *charonrpc.CreateUserRequest) *charonrpc.CreateUserResponse {
	res, err := suite.charon.user.Create(ctx, req)
	if err != nil {
		t.Fatalf("unexpected create user error: %s", err.Error())
	}
	if res.User.Id == 0 {
		t.Fatal("created user wrong id")
	} else {
		t.Logf("user has been created with id %d", res.User.Id)
	}

	return res
}

func testRPCServerCreateGroup(t *testing.T, suite *endToEndSuite, ctx context.Context, req *charonrpc.CreateGroupRequest) *charonrpc.CreateGroupResponse {
	res, err := suite.charon.group.Create(ctx, req)
	if err != nil {
		t.Fatalf("unexpected create group error: %s", err.Error())
	}
	if res.Group.Id == 0 {
		t.Fatal("created group wrong id")
	} else {
		t.Logf("group has been created with id %d", res.Group.Id)
	}

	return res
}

func testRPCServerRegisterPermissions(t *testing.T, suite *endToEndSuite, ctx context.Context, req *charonrpc.RegisterPermissionsRequest) *charonrpc.RegisterPermissionsResponse {
	res, err := suite.charon.permission.Register(ctx, req)
	if err != nil {
		t.Fatalf("unexpected permission registration error: %s", err.Error())
	}

	return res
}

func testRPCServerSetUserPermissions(t *testing.T, suite *endToEndSuite, ctx context.Context, req *charonrpc.SetUserPermissionsRequest) *charonrpc.SetUserPermissionsResponse {
	res, err := suite.charon.user.SetPermissions(ctx, req)
	if err != nil {
		t.Fatalf("unexpected set user permissions error: %s", err.Error())
	}

	return res
}

func testRPCServerSetUserGroups(t *testing.T, suite *endToEndSuite, ctx context.Context, req *charonrpc.SetUserGroupsRequest) *charonrpc.SetUserGroupsResponse {
	res, err := suite.charon.user.SetGroups(ctx, req)
	if err != nil {
		t.Fatalf("unexpected set user groups error: %s", err.Error())
	}

	return res
}
