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

var config configuration

func init() {
	config.init()
}

func main() {
	config.parse()

	logger := initLogger(config.logger.adapter, config.logger.format, config.logger.level, sklog.KeySubsystem, config.subsystem)
	postgres := initPostgres(
		config.postgres.connectionString,
		config.postgres.retry,
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

	userRepo := newUserRepository(postgres)
	permissionRepo := newPermissionRepository(postgres)

	// If any of this flags are set, try to create superuser. Will fail if data is wrong, or any user already exists.
	superuser := config.superuser
	if superuser.username != "" || superuser.password != "" || superuser.firstName != "" || superuser.lastName != "" {
		user, err := createSuperuser(userRepo, passwordHasher, superuser.username, superuser.password, superuser.firstName, superuser.lastName)
		if err != nil {
			sklog.Fatal(logger, err)
		}

		sklog.Info(logger, "superuser has been created", "username", user.Username, "first_name", user.FirstName, "last_name", user.LastName)
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
		logger:               logger,
		session:              mnemosyneClient,
		passwordHasher:       passwordHasher,
		userRepository:       userRepo,
		permissionRepository: permissionRepo,
	}
	charon.RegisterRPCServer(gRPCServer, charonServer)

	sklog.Info(logger, "rpc api is running", "host", config.host, "port", config.port, "subsystem", config.subsystem, "namespace", config.namespace)

	gRPCServer.Serve(listen)
}

func createSuperuser(userRepo *userRepository, hasher charon.PasswordHasher, username, plainPassword, firstName, lastName string) (*userEntity, error) {
	count, err := userRepo.Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("charond: superuser cannot be created, database is full of users")
	}

	securePassword, err := hasher.Hash(plainPassword)
	if err != nil {
		return nil, err
	}
	return userRepo.CreateSuperuser(username, securePassword, firstName, lastName)
}
