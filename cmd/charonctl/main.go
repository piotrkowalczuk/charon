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
	case "register":
		registerUser(config)
	default:
		fmt.Printf("unknown command %s", config.cmd())
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
	if !config.noauth {
		resp, err := client.Login(context.Background(), &charon.LoginRequest{
			Username: config.username,
			Password: config.password,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ctx = metadata.NewContext(ctx, metadata.Pairs(mnemosynerpc.AccessTokenMetadataKey, string(resp.AccessToken)))
	}

	return
}

func registerUser(config configuration) {
	c, ctx := client(config.address)
	resp, err := c.CreateUser(ctx, &charon.CreateUserRequest{
		Username:      config.register.username,
		PlainPassword: config.register.password,
		FirstName:     config.register.firstName,
		LastName:      config.register.lastName,
		IsSuperuser:   &ntypes.Bool{Bool: config.register.superuser, Valid: true},
	})
	if err != nil {
		fmt.Printf("registration failure: %s", grpc.ErrorDesc(err))
		os.Exit(1)
	}

	if config.register.superuser {
		fmt.Printf(`superuser "%s" has been created`, resp.User.Username)
	} else {
		fmt.Printf(`user "%s" has been created`, resp.User.Username)
	}
}
