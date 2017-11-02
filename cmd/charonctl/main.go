package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/ntypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
			// superuser does not need permissions
			os.Exit(0)
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

	list, err := c.group.List(ctx, &charonrpc.ListGroupsRequest{Limit: &ntypes.Int64{Valid: true}})
	if err != nil {
		return err
	}

FixturesLoop:
	for _, group := range fix.Groups {
		for _, existing := range list.Groups {
			if group.Name == existing.Name {
				// update permissions of already existing groups
				if err = setPermissions(ctx, c.group, existing.Id, existing.Name, group.Permissions); err != nil {
					return err
				}
				continue FixturesLoop
			}
		}

		res, err := c.group.Create(ctx, &charonrpc.CreateGroupRequest{
			Name:        group.Name,
			Description: &ntypes.String{Chars: group.Description, Valid: len(group.Description) > 0},
		})
		if err != nil {
			if grpc.Code(err) == codes.AlreadyExists {
				fmt.Printf("group (%s) already exists\n", group.Name)
				continue FixturesLoop
			}
			return fmt.Errorf("group (%s) creation failure: %s", group.Name, err.Error())
		}
		fmt.Printf("group (%s) has been created: %d\n", group.Name, res.Group.Id)
	}

	var perms []string
ExistingLoop:
	for _, existing := range list.Groups {
		for _, group := range fix.Groups {
			perms = group.Permissions
			if group.Name == existing.Name {
				continue ExistingLoop
			}
		}
		// set permissions to the groups that ware created
		if err = setPermissions(ctx, c.group, existing.Id, existing.Name, perms); err != nil {
			return err
		}
	}
	return nil
}

func setPermissions(ctx context.Context, group charonrpc.GroupManagerClient, id int64, name string, permissions []string) error {
	res, err := group.SetPermissions(ctx, &charonrpc.SetGroupPermissionsRequest{
		GroupId:     id,
		Permissions: permissions,
		Force:       true,
	})
	if err != nil {
		return fmt.Errorf("group (%s - %d) permission set failure: %s", name, id, err.Error())
	}
	fmt.Printf("group (%s) permissions has been set (given=%d,created=%d,removed=%d,untouched=%d)\n", name, len(permissions), res.Created, res.Removed, res.Untouched)

	got, err := group.ListPermissions(ctx, &charonrpc.ListGroupPermissionsRequest{Id: id})
	if err != nil {
		return fmt.Errorf("group (%s - %d) permission check failure: %s", name, id, err.Error())
	}
	var equal int
	for _, rp := range got.Permissions {
		for _, lp := range permissions {
			if rp == lp {
				equal++
			}
		}
	}
	if equal != len(permissions) {
		return fmt.Errorf("group (%s - %d) permission check failed, expected %d but got %d", name, id, len(permissions), equal)
	}
	return nil
}
