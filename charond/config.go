package main

import "flag"

// configuration ...
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
	monitoring struct {
		engine string
	}
	postgres struct {
		connectionString string
		retry            int
	}
}

// Init ...
func (c *configuration) Init() {
	if c == nil {
		*c = configuration{}
	}

	flag.StringVar(&c.host, "h", "127.0.0.1", "host")
	flag.IntVar(&c.port, "p", 8080, "port")
	flag.StringVar(&c.namespace, "n", "", "namespace")
	flag.StringVar(&c.subsystem, "s", "mnemosyne", "subsystem")
	flag.StringVar(&c.logger.adapter, "la", loggerAdapterStdOut, "logger adapter")
	flag.StringVar(&c.logger.format, "lf", loggerFormatJSON, "logger format")
	flag.IntVar(&c.logger.level, "ll", 6, "logger level")
	flag.StringVar(&c.monitoring.engine, "me", monitoringEnginePrometheus, "monitoring engine")
	flag.StringVar(&c.postgres.connectionString, "spcs", "postgres://localhost:5432?sslmode=disable", "storage postgres connection string")
	flag.IntVar(&c.postgres.retry, "spr", 10, "storage postgres possible attempts")
}

// Parse ...
func (c *configuration) Parse() {
	if !flag.Parsed() {
		flag.Parse()
	}
}
