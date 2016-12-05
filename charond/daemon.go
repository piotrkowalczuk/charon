package charond

import (
	"database/sql"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	libldap "github.com/go-ldap/ldap"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/ldap"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/promgrpc"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// DaemonOpts ...
type DaemonOpts struct {
	Test               bool
	Monitoring         bool
	TLS                bool
	TLSCertFile        string
	TLSKeyFile         string
	PostgresAddress    string
	PostgresDebug      bool
	PasswordBCryptCost int
	MnemosyneAddress   string
	LDAP               bool
	LDAPAddress        string
	LDAPBaseDN         string
	LDAPSearchDN       string
	LDAPBasePassword   string
	LDAPMappings       *ldap.Mappings
	Logger             log.Logger
	RPCListener        net.Listener
	DebugListener      net.Listener
}

// TestDaemonOpts represent set of options that can be passed to the TestDaemon constructor.
type TestDaemonOpts struct {
	MnemosyneAddress string
	PostgresAddress  string
}

// Daemon ...
type Daemon struct {
	opts          DaemonOpts
	monitor       *monitoring
	ldap          *libldap.Conn
	logger        log.Logger
	rpcListener   net.Listener
	debugListener net.Listener
	mnemosyneConn *grpc.ClientConn
	mnemosyne     mnemosynerpc.SessionManagerClient
}

// NewDaemon ...
func NewDaemon(opts DaemonOpts) *Daemon {
	d := &Daemon{
		opts:          opts,
		logger:        opts.Logger,
		rpcListener:   opts.RPCListener,
		debugListener: opts.DebugListener,
	}

	return d
}

// TestDaemon returns address of fully started in-memory daemon and closer to close it.
func TestDaemon(t *testing.T, opts TestDaemonOpts) (net.Addr, io.Closer) {
	l, err := net.Listen("tcp", "127.0.0.1:0") // any available address
	if err != nil {
		t.Fatalf("charon daemon tcp listener setup error: %s", err.Error())
	}

	logger := sklog.NewTestLogger(t)
	grpclog.SetLogger(sklog.NewGRPCLogger(logger))

	d := NewDaemon(DaemonOpts{
		Test:               true,
		Monitoring:         false,
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
	interceptor := promgrpc.NewInterceptor()
	if err = d.initMonitoring(); err != nil {
		return
	}

	var db *sql.DB
	db, err = initPostgres(d.opts.PostgresAddress, d.opts.Test, d.logger)
	if err != nil {
		return err
	}
	repos := newRepositories(db)

	d.mnemosyne, d.mnemosyneConn = initMnemosyne(d.opts.MnemosyneAddress, interceptor, d.logger)

	var passwordHasher password.Hasher
	if d.opts.LDAP {
		// open connection to check if it is reachable
		if d.ldap, err = initLDAP(d.opts.LDAPAddress, d.opts.LDAPBaseDN, d.opts.LDAPBasePassword, d.logger); err != nil {
			return
		}
		d.ldap.Close()
	}

	passwordHasher = initHasher(d.opts.PasswordBCryptCost, d.logger)
	if d.opts.Test {
		if _, err = createDummyTestUser(repos.user, passwordHasher); err != nil {
			return
		}
		sklog.Info(d.logger, "test super user has been created")
	}

	permissionReg := initPermissionRegistry(repos.permission, charon.AllPermissions, d.logger)

	opts := []grpc.ServerOption{
		// No stream endpoint available at the moment.
		grpc.UnaryInterceptor(unaryServerInterceptors(
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				start := time.Now()
				res, err := handler(ctx, req)

				if err != nil && grpc.Code(err) != codes.OK {
					sklog.Error(d.logger, errors.New(grpc.ErrorDesc(err)), "handler", info.FullMethod, "code", grpc.Code(err).String(), "elapsed", time.Since(start))
					return nil, handleError(err)
				}
				sklog.Debug(d.logger, "request handled successfully", "handler", info.FullMethod, "elapsed", time.Since(start))
				return res, err
			},
			interceptor.UnaryServer(),
		)),
	}
	if d.opts.TLS {
		creds, err := credentials.NewServerTLSFromFile(d.opts.TLSCertFile, d.opts.TLSKeyFile)
		if err != nil {
			return err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	grpclog.SetLogger(sklog.NewGRPCLogger(d.logger))
	gRPCServer := grpc.NewServer(opts...)
	server := &rpcServer{
		opts:               d.opts,
		logger:             d.logger,
		session:            d.mnemosyne,
		passwordHasher:     passwordHasher,
		permissionRegistry: permissionReg,
		repository:         repos,
		ldap: &sync.Pool{
			New: func() interface{} {
				conn, err := initLDAP(d.opts.LDAPAddress, d.opts.LDAPBaseDN, d.opts.LDAPBasePassword, d.logger)
				if err != nil {
					return err
				}
				return conn
			},
		},
	}
	charonrpc.RegisterAuthServer(gRPCServer, newAuth(server))
	charonrpc.RegisterUserManagerServer(gRPCServer, newUserManager(server))
	charonrpc.RegisterGroupManagerServer(gRPCServer, newGroupManager(server))
	charonrpc.RegisterPermissionManagerServer(gRPCServer, newPermissionManager(server))

	go func() {
		sklog.Info(d.logger, "rpc server is running", "address", d.rpcListener.Addr().String())

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
			sklog.Info(d.logger, "debug server is running", "address", d.debugListener.Addr().String())
			// TODO: implement keep alive

			mux := http.NewServeMux()
			mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
			mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
			mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
			mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
			mux.Handle("/metrics", prometheus.Handler())
			sklog.Error(d.logger, http.Serve(d.debugListener, mux))
		}()
	}

	return
}

// Close implements io.Closer interface.
func (d *Daemon) Close() (err error) {
	if d.ldap != nil {
		d.ldap.Close()
	}
	if err = d.mnemosyneConn.Close(); err != nil {
		return
	}
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
		return errors.New("getting hostname failed")
	}
	d.monitor = initPrometheus("charon", d.opts.Monitoring, prometheus.Labels{"server": hostname})
	return nil
}
