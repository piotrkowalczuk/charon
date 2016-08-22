package main

import (
	"fmt"
	"os"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/ntypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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

func client(addr string) (client charon.RPCClient, ctx context.Context) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client = charon.NewRPCClient(conn)
	ctx = context.Background()
	if !config.auth.disabled {
		resp, err := client.Login(context.Background(), &charon.LoginRequest{
			Username: config.auth.username,
			Password: config.auth.password,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ctx = metadata.NewContext(ctx, metadata.Pairs(mnemosynerpc.AccessTokenMetadataKey, resp.AccessToken))
	}

	return
}

func registerUser(config configuration) {
	c, ctx := client(config.address)
	res, err := c.CreateUser(ctx, &charon.CreateUserRequest{
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
		fmt.Printf("registration failure: %s\n", grpc.ErrorDesc(err))
		os.Exit(1)
	}

	if config.register.superuser {
		fmt.Printf("superuser has been created: %s\n", res.User.Username)
	} else {
		fmt.Printf("user has been created: %s\n", res.User.Username)
	}

	if config.register.superuser {
		resLogin, err := c.Login(ctx, &charon.LoginRequest{
			Username: config.register.username,
			Password: config.register.password,
			Client:   "charonctl",
		})
		if err != nil {
			fmt.Printf("login failure: %s\n", grpc.ErrorDesc(err))
			os.Exit(1)
		}
		ctx = metadata.NewContext(context.Background(), metadata.Pairs(mnemosynerpc.AccessTokenMetadataKey, resLogin.AccessToken))
	}

	if len(config.register.permissions) > 0 {
		_, err = c.SetUserPermissions(ctx, &charon.SetUserPermissionsRequest{
			UserId:      res.User.Id,
			Permissions: config.register.permissions.Strings(),
		})
		if err != nil {
			fmt.Printf("permission assigment failure: %s\n", grpc.ErrorDesc(err))
			os.Exit(1)
		}

		fmt.Println("users permissions has been set")
	}
}
