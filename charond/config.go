package main

import "flag"

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
		retry            int
	}
}

func (c *configuration) init() {
	if c == nil {
		*c = configuration{}
	}

	flag.StringVar(&c.host, "h", "127.0.0.1", "host")
	flag.IntVar(&c.port, "p", 8080, "port")
	flag.StringVar(&c.namespace, "n", "", "namespace")
	flag.StringVar(&c.subsystem, "s", "mnemosyne", "subsystem")
	flag.StringVar(&c.logger.adapter, "la", loggerAdapterStdOut, "logger adapter")
	flag.StringVar(&c.logger.format, "lf", loggerFormatJSON, "logger format")
	flag.StringVar(&c.mnemosyne.address, "ma", "", "mnemosyne session store connection address")
	flag.StringVar(&c.password.strategy, "ps", "bcrypt", "strategy how password will be stored")
	flag.IntVar(&c.password.bcrypt.cost, "pbc", 10, "bcrypt cost, bigget than safer (and longer to create)")
	flag.IntVar(&c.logger.level, "ll", 6, "logger level")
	flag.StringVar(&c.monitoring.engine, "me", monitoringEnginePrometheus, "monitoring engine")
	flag.StringVar(&c.postgres.connectionString, "pcs", "postgres://localhost:5432?sslmode=disable", "storage postgres connection string")
	flag.IntVar(&c.postgres.retry, "pr", 10, "storage postgres possible attempts")
}

func (c *configuration) parse() {
	if !flag.Parsed() {
		flag.Parse()
	}
}
