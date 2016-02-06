package main

import (
	"fmt"
	"log"

	"flag"

	"time"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
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

	rpc := charon.NewRPCClient(conn)
	if _, err = rpc.RegisterPermissions(context.Background(), &charon.RegisterPermissionsRequest{
		Permissions: []string{
			permissionCommentCanCreate.String(),
			permissionCommentCanEditAsOwner.String(),
			permissionCommentCanEditAsStranger.String(),
		},
	}); err != nil {
		log.Fatal(err)
	}

	sub, err := charon.New(conn).Subject(metadata.NewContext(context.Background(), metadata.Pairs("request_id", "123456789")), mnemosyne.DecodeToken([]byte(token)))
	if err != nil {
		log.Fatalf("%s: %s", grpc.Code(err).String(), grpc.ErrorDesc(err))
	}

	fmt.Printf("id: %d \n", sub.ID)
	fmt.Printf("username: %s \n", sub.Username)
	fmt.Printf("first name: %s \n", sub.FirstName)
	fmt.Printf("last name: %s \n", sub.LastName)
	fmt.Printf("is active: %t \n", sub.IsActive)
	fmt.Printf("is confirmed: %t \n", sub.IsConfirmed)
	fmt.Printf("is staff: %t \n", sub.IsStaff)
	fmt.Printf("is superuser: %t \n", sub.IsSuperuser)
	if len(sub.Permissions) > 0 {
		fmt.Println("permissions:")
	}
	for _, p := range sub.Permissions {
		fmt.Printf("     - %s \n", p.String())
	}
}
