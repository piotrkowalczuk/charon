package main

import (
	"fmt"
	"os"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/ntypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var config configuration

func init() {
	config.init()
}

func main() {
	config.parse()

	switch config.cmd() {
	case "help":
		config.cl.Usage()
	case "register":
		registerUser(config)
	default:
		fmt.Printf("unknown command %s\n", config.cmd())
	}
}

func registerUser(config configuration) {
	c, ctx := initClient(config.address)
	res, err := c.user.Create(ctx, &charonrpc.CreateUserRequest{
		Username:      config.register.username,
		PlainPassword: config.register.password,
		FirstName:     config.register.firstName,
		LastName:      config.register.lastName,
		IsSuperuser:   &ntypes.Bool{Bool: config.register.superuser, Valid: true},
		IsConfirmed:   &ntypes.Bool{Bool: config.register.confirmed, Valid: true},
		IsStaff:       &ntypes.Bool{Bool: config.register.staff, Valid: true},
		IsActive:      &ntypes.Bool{Bool: config.register.active, Valid: true},
	})
	if err != nil {
		if config.register.ifNotExists && grpc.Code(err) == codes.AlreadyExists {
			fmt.Printf("user already exists: %s\n", config.register.username)
			return
		}
		fmt.Printf("registration failure: %s\n", grpc.ErrorDesc(err))
		os.Exit(1)
	}

	if config.register.superuser {
		fmt.Printf("superuser has been created: %s\n", res.User.Username)
	} else {
		fmt.Printf("user has been created: %s\n", res.User.Username)
	}

	if len(config.register.permissions) > 0 {
		if config.register.superuser {
			token, err := c.auth.Login(ctx, &charonrpc.LoginRequest{
				Username: config.register.username,
				Password: config.register.password,
				Client:   "charonctl",
			})
			if err != nil {
				fmt.Printf("(superuser) login failure: %s\n", grpc.ErrorDesc(err))
				os.Exit(1)
			}
			ctx = metadata.NewContext(context.Background(), metadata.Pairs(mnemosyne.AccessTokenMetadataKey, token.Value))
		}

		if _, err = c.user.SetPermissions(ctx, &charonrpc.SetUserPermissionsRequest{
			UserId:      res.User.Id,
			Permissions: config.register.permissions.Strings(),
		}); err != nil {
			fmt.Printf("permission assigment failure: %s\n", grpc.ErrorDesc(err))
			os.Exit(1)
		}

		fmt.Println("users permissions has been set")
	}
}
