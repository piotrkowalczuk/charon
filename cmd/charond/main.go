package main

import (
	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/charon/charond"
	"github.com/piotrkowalczuk/sklog"
	"google.golang.org/grpc/grpclog"
)

var config configuration

func init() {
	config.init()
}

func main() {
	config.parse()

	logger := initLogger(config.logger.adapter, config.logger.format, config.logger.level)
	rpcListener := initListener(logger, config.host, config.port)
	debugListener := initListener(logger, config.host, config.port+1)

	grpclog.SetLogger(sklog.NewGRPCLogger(logger))

	daemon := charond.NewDaemon(charond.DaemonOpts{
		Test:                 config.test,
		TLS:                  config.tls.enabled,
		TLSCertFile:          config.tls.certFile,
		TLSKeyFile:           config.tls.keyFile,
		Monitoring:           config.monitoring.enabled,
		PostgresAddress:      config.postgres.address + "&application_name=charond_" + version,
		PostgresDebug:        config.postgres.debug,
		PasswordBCryptCost:   config.password.bcrypt.cost,
		MnemosyneAddress:     config.mnemosyned.address,
		MnemosyneTLS:         config.mnemosyned.tls.enabled,
		MnemosyneTLSCertFile: config.mnemosyned.tls.certFile,
		Logger:               logger,
		RPCListener:          rpcListener,
		DebugListener:        debugListener,
	})

	if err := daemon.Run(); err != nil {
		sklog.Fatal(logger, err)
	}
	defer daemon.Close()

	done := make(chan struct{})
	<-done
}
