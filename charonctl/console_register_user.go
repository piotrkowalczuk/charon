package charonctl

import (
	"context"
	"fmt"
	"os"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/ntypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type RegisterUserArg struct {
	IfNotExists bool
	Username    string
	Password    string
	FirstName   string
	LastName    string
	Superuser   bool
	Confirmed   bool
	Staff       bool
	Active      bool
	Permissions charon.Permissions
}

type consoleRegisterUser struct {
	user charonrpc.UserManagerClient
}

func (cru *consoleRegisterUser) RegisterUser(ctx context.Context, arg *RegisterUserArg) error {
	res, err := cru.user.Create(ctx, &charonrpc.CreateUserRequest{
		Username:      arg.Username,
		PlainPassword: arg.Password,
		FirstName:     arg.FirstName,
		LastName:      arg.LastName,
		IsSuperuser:   &ntypes.Bool{Bool: arg.Superuser, Valid: true},
		IsConfirmed:   &ntypes.Bool{Bool: arg.Confirmed, Valid: true},
		IsStaff:       &ntypes.Bool{Bool: arg.Staff, Valid: true},
		IsActive:      &ntypes.Bool{Bool: arg.Active, Valid: true},
	})
	if err != nil {
		if arg.IfNotExists && grpc.Code(err) == codes.AlreadyExists {
			return &Error{
				Msg: fmt.Sprintf("user already exists: %s\n", arg.Username),
				Err: err,
			}
		}
		return &Error{
			Msg: "registration failure",
			Err: err,
		}
	}

	if arg.Superuser {
		fmt.Printf("superuser has been created: %s\n", res.User.Username)
	} else {
		fmt.Printf("user has been created: %s\n", res.User.Username)
	}

	if len(arg.Permissions) > 0 {
		if arg.Superuser {
			// superuser does not need permissions
			os.Exit(0)
		}

		if _, err = cru.user.SetPermissions(ctx, &charonrpc.SetUserPermissionsRequest{
			UserId:      res.User.Id,
			Permissions: arg.Permissions.Strings(),
		}); err != nil {
			return &Error{
				Msg: "permission assigment failure",
				Err: err,
			}
		}

		fmt.Println("users permissions has been set")
	}
	return nil
}
