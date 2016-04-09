package main

import (
	"net"
	_ "net/http/pprof"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/sklog"
	"google.golang.org/grpc/grpclog"
)

var config configuration

func init() {
	config.init()
}

func main() {
	config.parse()

	logger := initLogger(config.logger.adapter, config.logger.format, config.logger.level, sklog.KeySubsystem, config.subsystem)
	rpcListener := initListener(logger, config.host, config.port)
	debugListener := initListener(logger, config.host, config.port+1)

	daemon := charon.NewDaemon(&charon.DaemonOpts{
		Namespace:          config.namespace,
		Subsystem:          config.subsystem,
		TLS:                config.tls.enabled,
		TLSCertFile:        config.tls.certFile,
		TLSKeyFile:         config.tls.keyFile,
		MonitoringEngine:   config.monitoring.engine,
		PostgresAddress:    config.postgres.address,
		PasswordBCryptCost: config.password.bcrypt.cost,
		MnemosyneAddress:   config.mnemosyne.address,
		Logger:             logger,
		RPCListener:        rpcListener,
		DebugListener:      debugListener,
	})

	grpclog.SetLogger(sklog.NewGRPCLogger(logger))
	if err := daemon.Run(); err != nil {
		sklog.Fatal(logger, err)
	}
	defer daemon.Close()

	done := make(chan struct{})
	<-done
}

func initListener(logger log.Logger, host string, port int) net.Listener {
	on := host + ":" + strconv.FormatInt(int64(port), 10)
	listener, err := net.Listen("tcp", on)
	if err != nil {
		sklog.Fatal(logger, err)
	}
	return listener
}
