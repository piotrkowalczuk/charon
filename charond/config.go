package main

import (
	"flag"
	"fmt"
	"os"
)

const VERSION="0.0.2"

type configuration struct {
	host      string
	port      int
	namespace string
	subsystem string
	logger    struct {
		adapter string
		format  string
		level   int
	}
	mnemosyne struct {
		address string
	}
	password struct {
		strategy string
		bcrypt   struct {
			cost int
		}
	}
	monitoring struct {
		engine string
	}
	postgres struct {
		connectionString string
	}
	superuser struct {
		username  string
		password  string
		firstName string
		lastName  string
	}
}

func (c *configuration) init() {
	if c == nil {
		*c = configuration{}
	}

	flag.StringVar(&c.host, "host", "127.0.0.1", "host")
	flag.IntVar(&c.port, "port", 8080, "port")
	flag.StringVar(&c.namespace, "namespace", "", "namespace")
	flag.StringVar(&c.subsystem, "subsystem", "charon", "subsystem")
	flag.StringVar(&c.logger.adapter, "l.adapter", loggerAdapterStdOut, "logger adapter")
	flag.StringVar(&c.logger.format, "l.format", loggerFormatJSON, "logger format")
	flag.IntVar(&c.logger.level, "l.level", 6, "logger level")
	flag.StringVar(&c.mnemosyne.address, "mnemo.address", "", "mnemosyne session store connection address")
	flag.StringVar(&c.password.strategy, "pwd.strategy", "bcrypt", "strategy how password will be stored")
	flag.IntVar(&c.password.bcrypt.cost, "pwd.bcryptcost", 10, "bcrypt cost, bigget than safer (and longer to create)")
	flag.StringVar(&c.monitoring.engine, "m.engine", monitoringEnginePrometheus, "monitoring engine")
	flag.StringVar(&c.postgres.connectionString, "ps.connectionstring", "postgres://localhost:5432?sslmode=disable", "storage postgres connection string")

	// Superuser configuration
	flag.StringVar(&c.superuser.username, "su.username", "", "superuser username")
	flag.StringVar(&c.superuser.password, "su.password", "", "superuser password")
	flag.StringVar(&c.superuser.firstName, "su.firstname", "", "superuser first name")
	flag.StringVar(&c.superuser.lastName, "su.lastname", "", "superuser last name")
}

func (c *configuration) parse() {

	if !flag.Parsed() {
		ver := flag.Bool("version", false, "Print version and exit")
		flag.Parse()
		if *ver {
			fmt.Printf("%s", VERSION)
			os.Exit(0)
		}
	}
}
