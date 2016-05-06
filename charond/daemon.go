package charond

import (
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const (
	// EnvironmentProduction ...
	EnvironmentProduction = "prod"
	// EnvironmentTest ...
	EnvironmentTest = "test"
)

// DaemonOpts ...
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

type TestDaemonOpts struct {
	MnemosyneAddress string
	PostgresAddress  string
}

// Daemon ...
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

// TestDaemon returns address of fully started in-memory daemon and closer to close it.
func TestDaemon(t *testing.T, opts *TestDaemonOpts) (net.Addr, io.Closer) {
	l, err := net.Listen("tcp", "127.0.0.1:0") // any available address
	if err != nil {
		t.Fatalf("charon daemon tcp listener setup error: %s", err.Error())
	}

	logger := sklog.NewTestLogger(t)
	grpclog.SetLogger(sklog.NewGRPCLogger(logger))

	d := NewDaemon(&DaemonOpts{
		Namespace:          "charon_test",
		Environment:        EnvironmentTest,
		MonitoringEngine:   MonitoringEnginePrometheus,
		MnemosyneAddress:   opts.MnemosyneAddress,
		Logger:             logger,
		PostgresAddress:    opts.PostgresAddress,
		RPCListener:        l,
		PasswordBCryptCost: bcrypt.MinCost,
	})
	if err := d.Run(); err != nil {
		t.Fatalf("charon daemon start error: %s", err.Error())
	}

	return d.Addr(), d
}

// Run ...
func (d *Daemon) Run() (err error) {
	var (
		mnemosyneClient mnemosyne.Mnemosyne
	)

	if err = d.initMonitoring(); err != nil {
		return
	}

	postgres, err := initPostgres(d.opts.PostgresAddress, d.opts.Environment, d.logger)
	if err != nil {
		return err
	}
	passwordHasher := initPasswordHasher(d.opts.PasswordBCryptCost, d.logger)
	d.mnemosyneConn, mnemosyneClient = initMnemosyne(d.opts.MnemosyneAddress, d.logger)

	repos := newRepositories(postgres)
	if d.opts.Environment == EnvironmentTest {
		if _, err = createDummyTestUser(repos.user, passwordHasher); err != nil {
			return
		}
		sklog.Info(d.logger, "test super user has been created")
	}

	permissionReg := initPermissionRegistry(repos.permission, charon.AllPermissions, d.logger)

	var opts []grpc.ServerOption
	if d.opts.TLS {
		creds, err := credentials.NewServerTLSFromFile(d.opts.TLSCertFile, d.opts.TLSKeyFile)
		if err != nil {
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpclog.SetLogger(sklog.NewGRPCLogger(d.logger))
	gRPCServer := grpc.NewServer(opts...)
	charonServer := &rpcServer{
		logger:             d.logger,
		session:            mnemosyneClient,
		passwordHasher:     passwordHasher,
		permissionRegistry: permissionReg,
		repository:         repos,
	}
	charon.RegisterRPCServer(gRPCServer, charonServer)

	go func() {
		sklog.Info(d.logger, "rpc server is running", "address", d.rpcListener.Addr().String(), "subsystem", d.opts.Subsystem, "namespace", d.opts.Namespace)

		if err := gRPCServer.Serve(d.rpcListener); err != nil {
			if err == grpc.ErrServerStopped {
				sklog.Info(d.logger, "grpc server has been stoped")
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
		return errors.New("charond: getting hostname failed")
	}

	switch d.opts.MonitoringEngine {
	case "":
		return errors.New("charond: monitoring is mandatory, at least for now")
	case MonitoringEnginePrometheus:
		d.monitor = initPrometheus(d.opts.Namespace, d.opts.Subsystem, prometheus.Labels{"server": hostname})
		return
	default:
		return errors.New("charond: unknown monitoring engine")
	}
}
