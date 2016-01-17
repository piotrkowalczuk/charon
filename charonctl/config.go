package main

import (
	"flag"
	"os"
)

type configuration struct {
	cl       *flag.FlagSet
	username string
	password string
	noauth   bool
	register struct {
		username  string
		password  string
		firstName string
		lastName  string
		superuser bool
	}
}

func (c *configuration) init() {
	*c = configuration{
		cl: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}

	c.cl.BoolVar(&c.noauth, "noauth", false, "noauth")
	c.cl.StringVar(&c.username, "username", "", "username")
	c.cl.StringVar(&c.password, "password", "", "password")
	c.cl.StringVar(&c.register.username, "r.username", "", "username")
	c.cl.StringVar(&c.register.password, "r.password", "", "password")
	c.cl.StringVar(&c.register.firstName, "r.firstname", "", "first name")
	c.cl.StringVar(&c.register.lastName, "r.lastname", "", "last name")
	c.cl.BoolVar(&c.register.superuser, "r.superuser", false, "superuser")
}

func (c *configuration) parse() {
	if c == nil || c.cl == nil {
		c.init()
	}
	if !c.cl.Parsed() {
		c.cl.Parse(os.Args[2:])
	}
}

func (c *configuration) cmd() string {
	return os.Args[1]
}
