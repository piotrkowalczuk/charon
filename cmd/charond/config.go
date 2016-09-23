package main

import (
	"flag"
	"fmt"
	"os"
)

const VERSION = "0.3.1"

type configuration struct {
	host      string
	port      int
	test      bool
	logger    struct {
		adapter string
		format  string
		level   int
	}
	mnemosyned struct {
		address string
	}
	password struct {
		strategy string
		bcrypt   struct {
			cost int
		}
	}
	monitoring struct {
		enabled bool
	}
	postgres struct {
		address string
		debug   bool
	}
	tls struct {
		enabled  bool
		certFile string
		keyFile  string
	}
	ldap struct {
		enabled  bool
		address  string
		dn       string
		password string
	}
}

func (c *configuration) init() {
	if c == nil {
		*c = configuration{}
	}

	flag.StringVar(&c.host, "host", "127.0.0.1", "host")
	flag.IntVar(&c.port, "port", 8080, "port")
	flag.BoolVar(&c.test, "test", false, "determines in what mode application starts")
	flag.StringVar(&c.logger.adapter, "log.adapter", loggerAdapterStdOut, "logger adapter")
	flag.StringVar(&c.logger.format, "log.format", loggerFormatJSON, "logger format")
	flag.IntVar(&c.logger.level, "log.level", 6, "logger level")
	flag.StringVar(&c.mnemosyned.address, "mnemosyned.address", "mnemosyned:8080", "mnemosyne daemon session store connection address")
	flag.StringVar(&c.password.strategy, "password.strategy", "bcrypt", "strategy how password will be stored")
	flag.IntVar(&c.password.bcrypt.cost, "password.bcryptcost", 10, "bcrypt cost, bigget than safer (and longer to create)")
	flag.BoolVar(&c.monitoring.enabled, "monitoring", false, "toggle application monitoring")
	flag.StringVar(&c.postgres.address, "postgres.address", "postgres://postgres:postgres@postgres/postgres?sslmode=disable", "postgres connection string")
	flag.BoolVar(&c.postgres.debug, "postgres.debug", false, "if true database queries are logged")
	flag.BoolVar(&c.tls.enabled, "tls", false, "tls enable flag")
	flag.StringVar(&c.tls.certFile, "tls.certfile", "", "path to tls cert file")
	flag.StringVar(&c.tls.keyFile, "tls.keyfile", "", "path to tls key file")
	flag.BoolVar(&c.ldap.enabled, "ldap", false, "ldap enable flag")
	flag.StringVar(&c.ldap.address, "ldap.address", "", "ldap server address")
	flag.StringVar(&c.ldap.dn, "ldap.dn", "", "ldap base distinguished name")
	flag.StringVar(&c.ldap.password, "ldap.password", "", "ldap password")
}

func (c *configuration) parse() {
	if !flag.Parsed() {
		ver := flag.Bool("version", false, "print version and exit")
		flag.Parse()
		if *ver {
			fmt.Printf("%s", VERSION)
			os.Exit(0)
		}
	}
}
