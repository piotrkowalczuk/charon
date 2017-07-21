package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestModifyUserHandler_Firewall_success(t *testing.T) {
	cases := []struct {
		hint   string
		req    charonrpc.ModifyUserRequest
		entity model.UserEntity
		user   model.UserEntity
		perm   charon.Permissions
	}{

		{
			hint:   "superuser should be able to degrade itself",
			req:    charonrpc.ModifyUserRequest{Id: 1, IsSuperuser: &ntypes.Bool{Bool: false, Valid: true}},
			entity: model.UserEntity{ID: 1, IsSuperuser: true},
			user:   model.UserEntity{ID: 1, IsSuperuser: true},
		},
		{
			hint:   "superuser should be able to promote an user",
			req:    charonrpc.ModifyUserRequest{Id: 2, IsSuperuser: &ntypes.Bool{Bool: true, Valid: true}},
			entity: model.UserEntity{ID: 2},
			user:   model.UserEntity{ID: 1, IsSuperuser: true},
		},
		{
			hint:   "superuser should be able to promote a staff user",
			req:    charonrpc.ModifyUserRequest{Id: 2, IsSuperuser: &ntypes.Bool{Bool: true, Valid: true}},
			entity: model.UserEntity{ID: 2},
			user:   model.UserEntity{ID: 1, IsSuperuser: true},
		},
		{
			hint:   "superuser should be able to degrade another superuser",
			req:    charonrpc.ModifyUserRequest{Id: 2, IsSuperuser: &ntypes.Bool{Bool: false, Valid: true}},
			entity: model.UserEntity{ID: 2, IsSuperuser: true},
			user:   model.UserEntity{ID: 1, IsSuperuser: true},
		},
		{
			hint:   "if user has permission to modify an user as a stranger, he should be able to do that",
			req:    charonrpc.ModifyUserRequest{Id: 2},
			entity: model.UserEntity{ID: 2},
			user:   model.UserEntity{ID: 1},
			perm:   charon.Permissions{charon.UserCanModifyAsStranger},
		},
		{
			hint:   "if user has permission to modify an user as an owner, he should be able to do that",
			req:    charonrpc.ModifyUserRequest{Id: 1},
			entity: model.UserEntity{ID: 1, CreatedBy: ntypes.Int64{Int64: 1, Valid: true}},
			user:   model.UserEntity{ID: 1, IsSuperuser: false},
			perm:   charon.Permissions{charon.UserCanModifyAsOwner},
		},
	}

	handler := &modifyUserHandler{}
	for _, args := range cases {
		t.Run(args.hint, func(t *testing.T) {
			msg, ok := handler.firewall(&args.req, &args.entity, &session.Actor{
				User:        &args.user,
				Permissions: args.perm,
			})
			if !ok {
				t.Errorf(msg)
			}
		})
	}
}

func TestModifyUserHandler_Firewall_failure(t *testing.T) {
	cases := []struct {
		hint   string
		req    charonrpc.ModifyUserRequest
		entity model.UserEntity
		user   model.UserEntity
		perm   charon.Permissions
		exp    string
	}{

		{
			hint:   "only superuser can modify superuser",
			req:    charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("new-first-name")},
			entity: model.UserEntity{ID: 1, IsSuperuser: true},
			user:   model.UserEntity{ID: 2},
			exp:    "only superuser can modify a superuser account",
		},
		{
			hint:   "to modify staff account as an owner permission is required",
			req:    charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("new-first-name")},
			entity: model.UserEntity{ID: 1, IsStaff: true},
			user:   model.UserEntity{ID: 1},
			perm: charon.Permissions{
				charon.UserCanModifyStaffAsStranger,
			},
			exp: "missing permission to modify staff account as an owner",
		},
		{
			hint:   "to modify staff account as a stranger permission is required",
			req:    charonrpc.ModifyUserRequest{Id: 1, FirstName: ntypes.NewString("new-first-name")},
			entity: model.UserEntity{ID: 1, IsStaff: true},
			user:   model.UserEntity{ID: 2},
			perm: charon.Permissions{
				charon.UserCanModifyStaffAsOwner,
			},
			exp: "missing permission to modify staff account as a stranger",
		},
		{
			hint:   "to modify staff account as a stranger permission is required",
			req:    charonrpc.ModifyUserRequest{Id: 1, IsSuperuser: ntypes.True()},
			entity: model.UserEntity{ID: 1},
			user:   model.UserEntity{ID: 2},
			exp:    "only superuser can change existing account to superuser",
		},
		{
			hint:   "to modify staff account as a stranger permission is required",
			req:    charonrpc.ModifyUserRequest{Id: 1, IsStaff: ntypes.True()},
			entity: model.UserEntity{ID: 1},
			user:   model.UserEntity{ID: 2},
			exp:    "missing permission to change existing account to staff",
		},
	}

	handler := &modifyUserHandler{}
	for _, args := range cases {
		t.Run(args.hint, func(t *testing.T) {
			msg, ok := handler.firewall(&args.req, &args.entity, &session.Actor{
				User:        &args.user,
				Permissions: args.perm,
			})
			if ok {
				t.Fatal("expected false")
			}
			if msg != args.exp {
				t.Errorf("wrong msg, expected '%s' bug to '%s'", args.exp, msg)
			}
		})
	}
}

func TestModifyUserHandler_Modify(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	cres := testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John",
		LastName:      "Snow",
	})
	_, err := suite.charon.user.Modify(ctx, &charonrpc.ModifyUserRequest{
		Id:        cres.User.Id,
		Username:  ntypes.NewString("john88@snow.com"),
		FirstName: ntypes.NewString("john"),
		LastName:  ntypes.NewString("snow"),
		IsActive:  ntypes.True(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestModifyUserHandler_Modify_nonExistingUser(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_ = testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John",
		LastName:      "Snow",
	})
	_, err := suite.charon.user.Modify(ctx, &charonrpc.ModifyUserRequest{
		Id:        1000,
		Username:  ntypes.NewString("john88@snow.com"),
		FirstName: ntypes.NewString("john"),
		LastName:  ntypes.NewString("snow"),
		IsActive:  ntypes.True(),
	})
	if err == nil {
		t.Fatal("missing error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.NotFound {
			t.Errorf("wrong error code, expected %s but got %s", codes.NotFound.String(), st.Code().String())
		}
	}
}

func TestModifyUserHandler_Modify_wrongID(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_ = testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John",
		LastName:      "Snow",
	})
	_, err := suite.charon.user.Modify(ctx, &charonrpc.ModifyUserRequest{
		Id: -1,
	})
	if err == nil {
		t.Fatal("missing error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.InvalidArgument {
			t.Errorf("wrong error code, expected %s but got %s", codes.InvalidArgument.String(), st.Code().String())
		}
	}
}

func TestModifyUserHandler_Modify_usernameAlreadyExists(t *testing.T) {
	suite := &endToEndSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	ctx := testRPCServerLogin(t, suite)

	_ = testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John",
		LastName:      "Snow",
	})
	cres := testRPCServerCreateUser(t, suite, timeout(ctx), &charonrpc.CreateUserRequest{
		Username:      "john2@snow.com",
		PlainPassword: "winteriscomming",
		FirstName:     "John2",
		LastName:      "Snow2",
	})
	_, err := suite.charon.user.Modify(ctx, &charonrpc.ModifyUserRequest{
		Id:       cres.User.Id,
		Username: ntypes.NewString("john@snow.com"),
	})
	if err == nil {
		t.Fatal("missing error")
	}
	if st, ok := status.FromError(err); ok {
		if st.Code() != codes.AlreadyExists {
			t.Errorf("wrong error code, expected %s but got %s", codes.AlreadyExists.String(), st.Code().String())
		}
	}
}
