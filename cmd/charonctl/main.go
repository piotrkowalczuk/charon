package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/grpc/status"

	"github.com/piotrkowalczuk/charon/charonctl"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/ntypes"
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
		ctl := connect(config)
		err := ctl.RegisterUser(ctl.Ctx, &charonctl.RegisterUserArg{
			Username:    config.register.username,
			Password:    config.register.password,
			IfNotExists: config.register.ifNotExists,
			Active:      config.register.active,
			Staff:       config.register.staff,
			Confirmed:   config.register.confirmed,
			Superuser:   config.register.superuser,
			FirstName:   config.register.firstName,
			LastName:    config.register.lastName,
			Permissions: config.register.permissions,
		})
		fail(err)
	case "refresh-token":
		ctl := connect(config)
		err := ctl.ObtainRefreshToken(ctl.Ctx, &charonctl.ObtainRefreshTokenArg{
			ExpireAfter: config.refreshToken.expireAfter,
			Notes:       config.refreshToken.notes,
		})
		fail(err)
	case "load":
		if err := load(config); err != nil {
			fmt.Printf("fixtures import failure: %s\n", status.Convert(err).Message())
			os.Exit(1)
		}
	default:
		fmt.Printf("unknown command %s\n", config.cmd())
	}
}

func fail(err error) {
	if err != nil {
		if ctlErr, ok := err.(*charonctl.Error); ok {
			fmt.Printf("%s: %s", ctlErr.Msg, status.Convert(ctlErr.Err).Message())
		} else {
			fmt.Println(status.Convert(err).Message())
		}
		os.Exit(1)
	}
}

func connect(config configuration) *charonctl.Console {
	conn, err := grpc.Dial(config.address, grpc.WithInsecure(), grpc.WithUserAgent("charonctl"))
	if err != nil {
		fmt.Printf("charond connection failure to %s with error: %s\n", config.address, status.Convert(err).Message())
		os.Exit(1)
	}

	var username, password string
	if config.auth.enabled {
		username = config.auth.username
		password = config.auth.password
	}

	c, err := charonctl.NewConsole(charonctl.ConsoleOpts{
		Conn:     conn,
		Username: username,
		Password: password,
	})
	fail(err)
	return c
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
			if status.Code(err) == codes.AlreadyExists {
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
