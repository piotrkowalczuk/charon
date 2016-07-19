package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/piotrkowalczuk/charon/charond"
)

const VERSION = "0.3.1"

type configuration struct {
	host      string
	port      int
	namespace string
	subsystem string
	test      bool
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
	flag.StringVar(&c.namespace, "namespace", "", "namespace")
	flag.StringVar(&c.subsystem, "subsystem", "charon", "subsystem")
	flag.BoolVar(&c.test, "test", false, "determines in what mode application starts")
	flag.StringVar(&c.logger.adapter, "log.adapter", loggerAdapterStdOut, "logger adapter")
	flag.StringVar(&c.logger.format, "log.format", loggerFormatJSON, "logger format")
	flag.IntVar(&c.logger.level, "log.level", 6, "logger level")
	flag.StringVar(&c.mnemosyne.address, "mnemosyne.address", "mnemosyne:8080", "mnemosyne session store connection address")
	flag.StringVar(&c.password.strategy, "password.strategy", "bcrypt", "strategy how password will be stored")
	flag.IntVar(&c.password.bcrypt.cost, "password.bcryptcost", 10, "bcrypt cost, bigget than safer (and longer to create)")
	flag.StringVar(&c.monitoring.engine, "monitoring.engine", charond.MonitoringEnginePrometheus, "monitoring engine")
	flag.StringVar(&c.postgres.address, "postgres.address", "postgres://postgres:postgres@postgres/postgres?sslmode=disable", "postgres connection string")
	flag.BoolVar(&c.postgres.debug, "postgres.debug", false, "if true database queries are logged")
	flag.BoolVar(&c.tls.enabled, "tls", false, "tls enable flag")
	flag.StringVar(&c.tls.certFile, "tls.certfile", "", "path to tls cert file")
	flag.StringVar(&c.tls.keyFile, "tls.keyfile", "", "path to tls key file")
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
