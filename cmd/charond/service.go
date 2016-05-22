package main

import (
	"io"
	"io/ioutil"
	stdlog "log"
	"net"
	"os"
	"strconv"

	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/sklog"
)

const (
	loggerAdapterStdOut = "stdout"
	loggerAdapterNone   = "none"
	loggerFormatJSON    = "json"
	loggerFormatHumane  = "humane"
	loggerFormatLogFmt  = "logfmt"
)

func initLogger(adapter, format string, level int, context ...interface{}) log.Logger {
	var (
		l log.Logger
		a io.Writer
	)

	switch adapter {
	case loggerAdapterStdOut:
		a = os.Stdout
	case loggerAdapterNone:
		a = ioutil.Discard
	default:
		stdlog.Fatal("charond: unsupported logger adapter")
	}

	switch format {
	case loggerFormatHumane:
		l = sklog.NewHumaneLogger(a, sklog.DefaultHTTPFormatter)
	case loggerFormatJSON:
		l = log.NewJSONLogger(a)
	case loggerFormatLogFmt:
		l = log.NewLogfmtLogger(a)
	default:
		stdlog.Fatal("charond: unsupported logger format")
	}

	l = log.NewContext(l).With(context...)

	sklog.Info(l, "logger has been initialized successfully", "adapter", adapter, "format", format, "level", level)

	return l
}

func initListener(logger log.Logger, host string, port int) net.Listener {
	on := host + ":" + strconv.FormatInt(int64(port), 10)
	listener, err := net.Listen("tcp", on)
	if err != nil {
		sklog.Fatal(logger, err)
	}
	return listener
}
