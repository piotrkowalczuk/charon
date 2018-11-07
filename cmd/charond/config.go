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
		environment string
		level       string
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
}

func (c *configuration) init() {
	if c == nil {
		*c = configuration{}
	}

	flag.StringVar(&c.host, "host", "127.0.0.1", "host")
	flag.IntVar(&c.port, "port", 8080, "port")
	flag.BoolVar(&c.test, "test", false, "determines in what mode application starts")
	// LOGGER
	flag.StringVar(&c.logger.environment, "log.environment", "production", "Logger environment config (production, stackdriver or development).")
	flag.StringVar(&c.logger.level, "log.level", "info", "Logger level (debug, info, warn, error, dpanic, panic, fatal)")
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
