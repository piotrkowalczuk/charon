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

	logger := initLogger(config.logger.adapter, config.logger.format, config.logger.level, sklog.KeySubsystem, config.subsystem)
	rpcListener := initListener(logger, config.host, config.port)
	debugListener := initListener(logger, config.host, config.port+1)

	daemon := charond.NewDaemon(charond.DaemonOpts{
		Test:                  config.test,
		TLS:                   config.tls.enabled,
		TLSCertFile:           config.tls.certFile,
		TLSKeyFile:            config.tls.keyFile,
		Monitoring:            config.monitoring.enabled,
		PostgresAddress:       config.postgres.address,
		PostgresDebug:         config.postgres.debug,
		PasswordBCryptCost:    config.password.bcrypt.cost,
		MnemosyneAddress:      config.mnemosyned.address,
		Logger:                logger,
		RPCListener:           rpcListener,
		DebugListener:         debugListener,
		LDAP:                  config.ldap.enabled,
		LDAPAddress:           config.ldap.address,
		LDAPDistinguishedName: config.ldap.dn,
		LDAPPassword:          config.ldap.password,
	})

	grpclog.SetLogger(sklog.NewGRPCLogger(logger))
	if err := daemon.Run(); err != nil {
		sklog.Fatal(logger, err)
	}
	defer daemon.Close()

	done := make(chan struct{})
	<-done
}
