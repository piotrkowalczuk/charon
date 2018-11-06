package main

import (
	"context"
	"fmt"
	"os"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/mnemosyne"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type client struct {
	auth       charonrpc.AuthClient
	user       charonrpc.UserManagerClient
	group      charonrpc.GroupManagerClient
	permission charonrpc.PermissionManagerClient
}

func initClient(addr string) (c *client, ctx context.Context) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithUserAgent("charonctl"))
	if err != nil {
		fmt.Printf("charond connection failure to %s with error: %s\n", addr, grpc.ErrorDesc(err))
		os.Exit(1)
	}

	c = &client{
		auth:       charonrpc.NewAuthClient(conn),
		user:       charonrpc.NewUserManagerClient(conn),
		group:      charonrpc.NewGroupManagerClient(conn),
		permission: charonrpc.NewPermissionManagerClient(conn),
	}
	ctx = context.Background()

	if config.auth.enabled {
		resp, err := c.auth.Login(context.Background(), &charonrpc.LoginRequest{
			Username: config.auth.username,
			Password: config.auth.password,
		})
		if err != nil {
			fmt.Printf("(initial) login failure: %s\n", grpc.ErrorDesc(err))
			os.Exit(1)
		}

		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(mnemosyne.AccessTokenMetadataKey, resp.Value))
	}

	return
}
