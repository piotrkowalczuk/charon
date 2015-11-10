package main

import (
	"errors"
	"net"
	"os"
	"strconv"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	config configuration
)

func init() {
	config.Init()
}

func main() {
	config.Parse()

	initLogger(config.logger.adapter, config.logger.format, config.logger.level, sklog.KeySubsystem, config.subsystem)
	initPostgres(
		config.postgres.connectionString,
		config.postgres.retry,
		logger,
	)

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

	listenOn := config.host + ":" + strconv.FormatInt(int64(config.port), 10)
	listen, err := net.Listen("tcp", listenOn)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	var opts []grpc.ServerOption
	//	if *tls {
	//		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	//		if err != nil {
	//			grpclog.Fatalf("Failed to generate credentials %v", err)
	//		}
	//		opts = []grpc.ServerOption{grpc.Creds(creds)}
	//	}
	grpclog.SetLogger(sklog.NewGRPCLogger(logger))
	gRPCServer := grpc.NewServer(opts...)

	charonServer := &rpcServer{
		logger: logger,
	}
	charon.RegisterRPCServer(gRPCServer, charonServer)

	gRPCServer.Serve(listen)
}
