package charon

import (
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	EnvironmentProduction = "prod"
	EnvironmentTest       = "test"
)

type DaemonOpts struct {
	Environment        string
	Namespace          string
	Subsystem          string
	MonitoringEngine   string
	TLS                bool
	TLSCertFile        string
	TLSKeyFile         string
	PostgresAddress    string
	PasswordBCryptCost int
	MnemosyneAddress   string
	Logger             log.Logger
	RPCListener        net.Listener
	DebugListener      net.Listener
}

type Daemon struct {
	opts          *DaemonOpts
	monitor       *monitoring
	logger        log.Logger
	rpcListener   net.Listener
	debugListener net.Listener
	mnemosyneConn *grpc.ClientConn
}

// NewDaemon ...
func NewDaemon(opts *DaemonOpts) *Daemon {
	d := &Daemon{
		opts:          opts,
		logger:        opts.Logger,
		rpcListener:   opts.RPCListener,
		debugListener: opts.DebugListener,
	}

	return d
}

func (d *Daemon) Run() (err error) {
	var (
		mnemosyneClient mnemosyne.Mnemosyne
	)

	if err = d.initMonitoring(); err != nil {
		return
	}

	postgres, err := initPostgres(d.opts.PostgresAddress, d.opts.Environment, d.logger)
	if err != nil {
		sklog.Error(d.logger, err)
	}
	passwordHasher := initPasswordHasher(d.opts.PasswordBCryptCost, d.logger)
	d.mnemosyneConn, mnemosyneClient = initMnemosyne(d.opts.MnemosyneAddress, d.logger)

	repos := newRepositories(postgres)
	if d.opts.Environment == EnvironmentTest {
		if _, err = createDumyTestUser(repos.user, passwordHasher); err != nil {
			return
		}
		sklog.Info(d.logger, "test super user has been created")
	}

	permissionReg := initPermissionRegistry(repos.permission, AllPermissions, d.logger)

	var opts []grpc.ServerOption
	if d.opts.TLS {
		creds, err := credentials.NewServerTLSFromFile(d.opts.TLSCertFile, d.opts.TLSKeyFile)
		if err != nil {
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	gRPCServer := grpc.NewServer(opts...)

	charonServer := &rpcServer{
		logger:             d.logger,
		session:            mnemosyneClient,
		passwordHasher:     passwordHasher,
		permissionRegistry: permissionReg,
		repository:         repos,
	}
	RegisterRPCServer(gRPCServer, charonServer)

	go func() {
		sklog.Info(d.logger, "rpc server is running", "address", d.rpcListener.Addr().String(), "subsystem", d.opts.Subsystem, "namespace", d.opts.Namespace)

		if err := gRPCServer.Serve(d.rpcListener); err != nil {
			if err == grpc.ErrServerStopped {
				return
			}

			sklog.Error(d.logger, err)
		}
	}()

	if d.debugListener != nil {
		go func() {
			sklog.Info(d.logger, "debug server is running", "address", d.debugListener.Addr().String(), "subsystem", d.opts.Subsystem, "namespace", d.opts.Namespace)
			// TODO: implement keep alive
			sklog.Error(d.logger, http.Serve(d.debugListener, nil))
		}()
	}

	return
}

// Close implements io.Closer interface.
func (d *Daemon) Close() (err error) {
	if err = d.rpcListener.Close(); err != nil {
		return
	}
	if d.debugListener != nil {
		err = d.debugListener.Close()
	}
	return
}

// Addr returns net.Addr that rpc service is listening on.
func (d *Daemon) Addr() net.Addr {
	return d.rpcListener.Addr()
}

func (d *Daemon) initMonitoring() (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return errors.New("charon: getting hostname failed")
	}

	switch d.opts.MonitoringEngine {
	case "":
		return errors.New("charon: monitoring is mandatory, at least for now")
	case MonitoringEnginePrometheus:
		d.monitor = initPrometheus(d.opts.Namespace, d.opts.Subsystem, prometheus.Labels{"server": hostname})
		return
	default:
		return errors.New("charon: unknown monitoring engine")
	}
}
