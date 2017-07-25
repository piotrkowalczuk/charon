package charontest_test

import (
	"testing"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/charontest"
)

func TestAuthClient(t *testing.T) {
	var mock interface{} = &charontest.AuthClient{}

	if _, ok := mock.(charonrpc.AuthClient); !ok {
		t.Error("auth client mock should implement original interface, but does not")
	}
}

func TestUserManagerClient(t *testing.T) {
	var mock interface{} = &charontest.UserManagerClient{}

	if _, ok := mock.(charonrpc.UserManagerClient); !ok {
		t.Error("user manager client mock should implement original interface, but does not")
	}
}

func TestGroupManagerClient(t *testing.T) {
	var mock interface{} = &charontest.GroupManagerClient{}

	if _, ok := mock.(charonrpc.GroupManagerClient); !ok {
		t.Error("group manager client mock should implement original interface, but does not")
	}
}

func TestPermissionManagerClient(t *testing.T) {
	var mock interface{} = &charontest.PermissionManagerClient{}

	if _, ok := mock.(charonrpc.PermissionManagerClient); !ok {
		t.Error("permission manager client mock should implement original interface, but does not")
	}
}
