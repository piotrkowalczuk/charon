package main

import (
	stdlog "log"

	"os"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
)

const (
	loggerAdapterStdOut = "stdout"
	loggerFormatJSON    = "json"
	loggerFormatHumane  = "humane"
	loggerFormatLogFmt  = "logfmt"
)

var logger log.Logger

func initLogger(adapter, format string, level int, context ...interface{}) log.Logger {
	var l log.Logger

	if adapter != loggerAdapterStdOut {
		stdlog.Fatal("service: unsupported logger adapter")
	}

	switch format {
	case loggerFormatHumane:
		l = sklog.NewHumaneLogger(os.Stdout)
	case loggerFormatJSON:
		l = log.NewJSONLogger(os.Stdout)
	case loggerFormatLogFmt:
		l = log.NewLogfmtLogger(os.Stdout)
	default:
		stdlog.Fatal("charond: unsupported logger format")
	}

	l = log.NewContext(l).With(context...)

	sklog.Info(l, "logger has been initialized successfully", "adapter", adapter, "format", format, "level", level)

	return l
}
