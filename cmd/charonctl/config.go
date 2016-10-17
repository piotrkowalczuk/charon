package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/piotrkowalczuk/charon"
)

type configuration struct {
	cl      *flag.FlagSet
	address string
	auth    struct {
		username string
		password string
		enabled  bool
	}
	register struct {
		ifNotExists bool
		username    string
		password    string
		firstName   string
		lastName    string
		superuser   bool
		confirmed   bool
		staff       bool
		active      bool
		permissions charon.Permissions
	}
}

func (c *configuration) init() {
	*c = configuration{
		cl: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}

	c.cl.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		c.cl.PrintDefaults()
	}
	c.cl.StringVar(&c.address, "address", "charond:8080", "charon address")
	c.cl.BoolVar(&c.auth.enabled, "auth", true, "authorization check flag")
	c.cl.StringVar(&c.auth.username, "auth.username", "", "username")
	c.cl.StringVar(&c.auth.password, "auth.password", "", "password")
	c.cl.BoolVar(&c.register.ifNotExists, "register.ifnotexists", false, "application does not fail if user already exists")
	c.cl.StringVar(&c.register.username, "register.username", "", "username")
	c.cl.StringVar(&c.register.password, "register.password", "", "password")
	c.cl.StringVar(&c.register.firstName, "register.firstname", "", "first name")
	c.cl.StringVar(&c.register.lastName, "register.lastname", "", "last name")
	c.cl.Var(&c.register.permissions, "register.permission", "list of permissions that user should")
	c.cl.BoolVar(&c.register.superuser, "register.superuser", false, "is user the superuser")
	c.cl.BoolVar(&c.register.confirmed, "register.confirmed", false, "is user account confirmed")
	c.cl.BoolVar(&c.register.staff, "register.staff", false, "is user part of the staff")
	c.cl.BoolVar(&c.register.active, "register.active", false, "is user account active")
}

func (c *configuration) parse() {
	if c == nil || c.cl == nil {
		c.init()
	}
	if !c.cl.Parsed() {
		if len(os.Args) > 1 {
			c.cl.Parse(os.Args[2:])
		}
	}
}

func (c *configuration) cmd() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return "help"
}
