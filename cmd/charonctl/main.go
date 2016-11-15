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
	"gopkg.in/square/go-jose.v1/json"
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
	case "load":
		if err := load(config); err != nil {
			fmt.Printf("fixtures import failure: %s\n", grpc.ErrorDesc(err))
			os.Exit(1)
		}
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

type fixtures struct {
	Groups []struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Permissions []string `json:"permissions"`
	} `json:"groups"`
}

func load(config configuration) error {
	c, ctx := initClient(config.address)

	file, err := os.Open(config.fixtures.path)
	if err != nil {
		return err
	}

	fix := fixtures{}
	if err = json.NewDecoder(file).Decode(&fix); err != nil {
		return err
	}

	for _, group := range fix.Groups {
		res, err := c.group.Create(ctx, &charonrpc.CreateGroupRequest{
			Name:        group.Name,
			Description: &ntypes.String{String: group.Description, Valid: len(group.Description) > 0},
		})
		if err != nil {
			if grpc.Code(err) == codes.AlreadyExists {
				fmt.Printf("group (%s) already exists\n", group.Name)
				continue
			}
			return fmt.Errorf("group (%s) creation failure: %s", group.Name, err.Error())
		}
		fmt.Printf("group (%s) has been created: %d\n", group.Name, res.Group.Id)

		_, err = c.group.SetPermissions(ctx, &charonrpc.SetGroupPermissionsRequest{
			GroupId:     res.Group.Id,
			Permissions: group.Permissions,
		})
		if err != nil {
			return fmt.Errorf("group (%s - %d) permission set failure: %s", group.Name, res.Group.Id, err.Error())
		}
		fmt.Println("group permissions has been set")
	}
	return nil
}
