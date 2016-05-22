package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestRPCServer_minimal(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServer_login(t, suite)
	permissions := []string{
		"winterfell:castle:can enter as a lord",
		"winterfell:castle:can close as a lord",
	}

	createUserResponse := testRPCServer_createUser(t, suite, ctx, &charon.CreateUserRequest{
		Username:      "john@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John",
		LastName:      "Snow",
	})
	createGroupResponse := testRPCServer_createGroup(t, suite, ctx, &charon.CreateGroupRequest{
		Name: "winterfell",
	})
	registerPermissionsResponse := testRPCServer_registerPermissions(t, suite, ctx, &charon.RegisterPermissionsRequest{
		Permissions: permissions,
	})
	if registerPermissionsResponse.Created != 2 {
		t.Fatalf("wrong number of registered permissions, expected 2 but got %d", registerPermissionsResponse.Created)
	}
	_ = testRPCServer_setUserPermissions(t, suite, ctx, &charon.SetUserPermissionsRequest{
		UserId:      createUserResponse.User.Id,
		Permissions: permissions,
	})
	_ = testRPCServer_setUserGroups(t, suite, ctx, &charon.SetUserGroupsRequest{
		UserId: createUserResponse.User.Id,
		Groups: []int64{createGroupResponse.Group.Id},
	})
}

func testRPCServer_login(t *testing.T, suite *endToEndSuite) context.Context {
	res, err := suite.charon.Login(context.TODO(), &charon.LoginRequest{Username: "test", Password: "test"})
	if err != nil {
		t.Fatalf("unexpected login error: %s: with code %s", grpc.ErrorDesc(err), grpc.Code(err))
	}
	meta := metadata.Pairs(mnemosynerpc.AccessTokenMetadataKey, res.AccessToken.Encode())
	return metadata.NewContext(context.Background(), meta)
}

func testRPCServer_createUser(t *testing.T, suite *endToEndSuite, ctx context.Context, req *charon.CreateUserRequest) *charon.CreateUserResponse {
	res, err := suite.charon.CreateUser(ctx, req)
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

func testRPCServer_createGroup(t *testing.T, suite *endToEndSuite, ctx context.Context, req *charon.CreateGroupRequest) *charon.CreateGroupResponse {
	res, err := suite.charon.CreateGroup(ctx, req)
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

func testRPCServer_registerPermissions(t *testing.T, suite *endToEndSuite, ctx context.Context, req *charon.RegisterPermissionsRequest) *charon.RegisterPermissionsResponse {
	res, err := suite.charon.RegisterPermissions(ctx, req)
	if err != nil {
		t.Fatalf("unexpected permission registration error: %s", err.Error())
	}

	return res
}

func testRPCServer_setUserPermissions(t *testing.T, suite *endToEndSuite, ctx context.Context, req *charon.SetUserPermissionsRequest) *charon.SetUserPermissionsResponse {
	res, err := suite.charon.SetUserPermissions(ctx, req)
	if err != nil {
		t.Fatalf("unexpected set user permissions error: %s", err.Error())
	}

	return res
}

func testRPCServer_setUserGroups(t *testing.T, suite *endToEndSuite, ctx context.Context, req *charon.SetUserGroupsRequest) *charon.SetUserGroupsResponse {
	res, err := suite.charon.SetUserGroups(ctx, req)
	if err != nil {
		t.Fatalf("unexpected set user groups error: %s", err.Error())
	}

	return res
}
