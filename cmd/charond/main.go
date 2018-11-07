package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/charon/internal/charond"
	"github.com/piotrkowalczuk/charon/internal/service/logger"
	"go.uber.org/zap"
)

var (
	config  configuration
	service = "charond"
)

func init() {
	config.init()
}

func main() {
	config.parse()

	log, err := logger.Init(service, version, logger.Opts{
		Environment: config.logger.environment,
		Level:       config.logger.level,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rpcListener := initListener(log, config.host, config.port)
	debugListener := initListener(log, config.host, config.port+1)

	// TODO: update and make it optional
	//grpclog.SetLogger(sklog.NewGRPCLogger(logger))

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
		Logger:               log.Named("daemon"),
		RPCListener:          rpcListener,
		DebugListener:        debugListener,
	})

	if err := daemon.Run(); err != nil {
		log.Fatal("daemon returned an error", zap.Error(err))
	}
	defer daemon.Close()

	done := make(chan struct{})
	<-done
}
