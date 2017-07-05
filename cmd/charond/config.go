package main

import (
	"flag"
	"fmt"
	"os"
)

var version string

type configuration struct {
	host   string
	port   int
	test   bool
	logger struct {
		adapter string
		format  string
		level   int
	}
	mnemosyned struct {
		address string
		tls     struct {
			enabled  bool
			certFile string
		}
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
		enabled bool
		address string
		search  string
		base    struct {
			dn, password string
		}
		mappings string
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
	// MNEMOSYNE
	flag.StringVar(&c.mnemosyned.address, "mnemosyned.address", "mnemosyned:8080", "mnemosyne daemon session store connection address")
	flag.BoolVar(&c.mnemosyned.tls.enabled, "mnemosyned.tls", false, "tls enable flag for mnemosyned client connection")
	flag.StringVar(&c.mnemosyned.tls.certFile, "mnemosyned.tls.crt", "", "path to tls cert file for mnemosyned client connection")
	// PASSWORD
	flag.StringVar(&c.password.strategy, "password.strategy", "bcrypt", "strategy how password will be stored")
	flag.IntVar(&c.password.bcrypt.cost, "password.bcryptcost", 10, "bcrypt cost, bigget than safer (and longer to create)")
	flag.BoolVar(&c.monitoring.enabled, "monitoring", false, "toggle application monitoring")
	// POSTGRES
	flag.StringVar(&c.postgres.address, "postgres.address", "postgres://postgres:postgres@postgres/postgres?sslmode=disable", "postgres connection string")
	flag.BoolVar(&c.postgres.debug, "postgres.debug", false, "if true database queries are logged")
	// TLS
	flag.BoolVar(&c.tls.enabled, "tls", false, "tls enable flag")
	flag.StringVar(&c.tls.certFile, "tls.crt", "", "path to tls cert file")
	flag.StringVar(&c.tls.keyFile, "tls.key", "", "path to tls key file")
	// LDAP
	flag.BoolVar(&c.ldap.enabled, "ldap", false, "ldap enable flag")
	flag.StringVar(&c.ldap.address, "ldap.address", "", "ldap server address")
	flag.StringVar(&c.ldap.base.dn, "ldap.base.dn", "", "ldap base distinguished name")
	flag.StringVar(&c.ldap.base.password, "ldap.base.password", "", "ldap password")
	flag.StringVar(&c.ldap.search, "ldap.search", "", "ldap search distinguished name")
	flag.StringVar(&c.ldap.mappings, "ldap.mappings", "", "path to the ldap mappings")
}

func (c *configuration) parse() {
	if !flag.Parsed() {
		ver := flag.Bool("version", false, "print version and exit")
		flag.Parse()
		if *ver {
			fmt.Printf("%s\n", version)
			os.Exit(0)
		}
	}
}
