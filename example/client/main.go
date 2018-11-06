package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	permissionCommentCanCreate         charon.Permission = "forumservice:comment:can create"
	permissionCommentCanEditAsOwner    charon.Permission = "forumservice:comment:can edit as an owner"
	permissionCommentCanEditAsStranger charon.Permission = "forumservice:comment:can edit as a stranger"
)

var (
	address string
	token   string
)

func init() {
	flag.StringVar(&address, "address", "localhost:8010", "charond service address")
	flag.StringVar(&token, "token", "", "session token")
}

func main() {
	flag.Parse()

	if token == "" {
		log.Fatal("missing sesion token")
	}

	conn, err := grpc.Dial(address, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		log.Fatal(grpc.ErrorDesc(err))
	}
	defer conn.Close()

	rpc := charonrpc.NewPermissionManagerClient(conn)
	if _, err = rpc.Register(context.Background(), &charonrpc.RegisterPermissionsRequest{
		Permissions: []string{
			permissionCommentCanCreate.String(),
			permissionCommentCanEditAsOwner.String(),
			permissionCommentCanEditAsStranger.String(),
		},
	}); err != nil {
		log.Fatal(err)
	}

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("request_id", "123456789"))
	res, err := charonrpc.NewAuthClient(conn).Actor(ctx, &wrappers.StringValue{Value: token})
	if err != nil {
		log.Fatalf("%s: %s", grpc.Code(err).String(), grpc.ErrorDesc(err))
	}

	fmt.Printf("id: %d \n", res.Id)
	fmt.Printf("username: %s \n", res.Username)
	fmt.Printf("first name: %s \n", res.FirstName)
	fmt.Printf("last name: %s \n", res.LastName)
	fmt.Printf("is active: %t \n", res.IsActive)
	fmt.Printf("is confirmed: %t \n", res.IsConfirmed)
	fmt.Printf("is staff: %t \n", res.IsStuff)
	fmt.Printf("is superuser: %t \n", res.IsSuperuser)
	if len(res.Permissions) > 0 {
		fmt.Println("permissions:")
	}
	for _, p := range res.Permissions {
		fmt.Printf("     - %s \n", p)
	}
}
