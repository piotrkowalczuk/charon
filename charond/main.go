package main

//go:generate charong
//go:generate mockery -all -inpkg -output_file=mocks_test.go

import (
	"errors"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

var config configuration

func init() {
	config.init()
}

func main() {
	config.parse()

	logger := initLogger(config.logger.adapter, config.logger.format, config.logger.level, sklog.KeySubsystem, config.subsystem)
	postgres := initPostgres(
		config.postgres.connectionString,
		config.test,
		logger,
	)
	passwordHasher := initPasswordHasher(config.password.bcrypt.cost, logger)
	mnemosyneConn, mnemosyneClient := initMnemosyne(config.mnemosyne.address, logger)

	defer mnemosyneConn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		sklog.Fatal(logger, errors.New("charond: getting hostname failed"))
	}

	switch config.monitoring.engine {
	case "":
		sklog.Fatal(logger, errors.New("charond: monitoring is mandatory, at least for now..."))
	case monitoringEnginePrometheus:
		initMonitoring(initPrometheus(config.namespace, config.subsystem, prometheus.Labels{"server": hostname}), logger)
	default:
		sklog.Fatal(logger, errors.New("charond: unknown monitoring engine"))
	}

	repos := newRepositories(postgres)
	if config.test {
		if _, err = createDumyTestUser(repos.user, passwordHasher); err != nil {
			sklog.Fatal(logger, err)
		}
		sklog.Info(logger, "test super user ")
	}

	permissionReg := initPermissionRegistry(repos.permission, charon.AllPermissions, logger)

	listenOn := config.host + ":" + strconv.FormatInt(int64(config.port), 10)
	listen, err := net.Listen("tcp", listenOn)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	var opts []grpc.ServerOption
	if config.tls.enabled {
		creds, err := credentials.NewServerTLSFromFile(config.tls.certFile, config.tls.keyFile)
		if err != nil {
			sklog.Fatal(logger, err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	//opts = append(opts, grpc.Creds(charon.NewCredentials()))
	grpclog.SetLogger(sklog.NewGRPCLogger(logger))
	gRPCServer := grpc.NewServer(opts...)

	charonServer := &rpcServer{
		logger:             logger,
		session:            mnemosyneClient,
		passwordHasher:     passwordHasher,
		permissionRegistry: permissionReg,
		repository:         repos,
	}
	charon.RegisterRPCServer(gRPCServer, charonServer)

	sklog.Info(logger, "rpc api is running", "host", config.host, "port", config.port, "subsystem", config.subsystem, "namespace", config.namespace)

	go func() {
		sklog.Fatal(logger, http.ListenAndServe(address(config.host, config.port+1), nil))
	}()
	gRPCServer.Serve(listen)
}
