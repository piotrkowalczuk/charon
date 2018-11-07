package charond

import (
	"context"
	"database/sql"
	"io"
	"net"
	"net/http"
	"net/http/pprof"
	"testing"
	"time"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/promgrpc/v3"
	"github.com/piotrkowalczuk/zapstackdriver/zapstackdrivergrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// DaemonOpts ...
type DaemonOpts struct {
	Test                 bool
	Monitoring           bool
	TLS                  bool
	TLSCertFile          string
	TLSKeyFile           string
	PostgresAddress      string
	PostgresDebug        bool
	PasswordBCryptCost   int
	MnemosyneAddress     string
	MnemosyneTLS         bool
	MnemosyneTLSCertFile string
	Logger               *zap.Logger
	RPCListener          net.Listener
	DebugListener        net.Listener
}

// TestDaemonOpts represent set of options that can be passed to the TestDaemon constructor.
type TestDaemonOpts struct {
	MnemosyneAddress string
	PostgresAddress  string
	PostgresDebug    bool
}

// Daemon ...
type Daemon struct {
	opts          DaemonOpts
	logger        *zap.Logger
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

	d := NewDaemon(DaemonOpts{
		Test:               true,
		Monitoring:         false,
		MnemosyneAddress:   opts.MnemosyneAddress,
		Logger:             zap.L(), // TODO: implement properly
		PostgresAddress:    opts.PostgresAddress,
		PostgresDebug:      opts.PostgresDebug,
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
	interceptor := promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})

	clientOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithStatsHandler(interceptor),
		grpc.WithTimeout(10 * time.Second),
		grpc.WithUserAgent("charond"),
		grpc.WithDialer(interceptor.Dialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("tcp", addr, timeout)
		})),
		grpc.WithUnaryInterceptor(interceptor.UnaryClient()),
		grpc.WithStreamInterceptor(interceptor.StreamClient()),
	}
	serverOpts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(100),
		grpc.StatsHandler(interceptor),
		// No stream endpoint available at the moment.
		grpc.UnaryInterceptor(unaryServerInterceptors(
			grpcerr.UnaryServerInterceptor(),
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				if md, ok := metadata.FromIncomingContext(ctx); ok {
					ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
						mnemosyne.AccessTokenMetadataKey: md[mnemosyne.AccessTokenMetadataKey],
						"request_id":                     md["request_id"],
					})
				}

				return handler(ctx, req)
			},
			zapstackdrivergrpc.UnaryServerInterceptor(d.logger),
			interceptor.UnaryServer(),
		)),
	}
	if d.opts.TLS {
		serverCreds, err := credentials.NewServerTLSFromFile(d.opts.TLSCertFile, d.opts.TLSKeyFile)
		if err != nil {
			return err
		}
		serverOpts = append(serverOpts, grpc.Creds(serverCreds))
	}

	if d.opts.MnemosyneTLS {
		var serverNameOverride string
		if d.opts.Test {
			serverNameOverride = "test.local.tld"
		}
		clientCreds, err := credentials.NewClientTLSFromFile(d.opts.MnemosyneTLSCertFile, serverNameOverride)
		if err != nil {
			return err
		}
		clientOpts = append(clientOpts, grpc.WithTransportCredentials(clientCreds))
	} else {
		clientOpts = append(clientOpts, grpc.WithInsecure())
	}

	var db *sql.DB
	db, err = initPostgres(d.opts.PostgresAddress, d.opts.Test, d.logger)
	if err != nil {
		return err
	}
	repos := newRepositories(db)

	d.mnemosyne, d.mnemosyneConn = initMnemosyne(d.opts.MnemosyneAddress, d.logger, clientOpts)

	passwordHasher := initHasher(d.opts.PasswordBCryptCost, d.logger)
	if d.opts.Test {
		if err = createDummyTestUser(
			context.TODO(),
			repos.user,
			repos.refreshToken,
			passwordHasher,
		); err != nil {
			return
		}
		d.logger.Info("test super user has been created")
	}

	permissionReg := initPermissionRegistry(repos.permission, charon.AllPermissions, d.logger)

	gRPCServer := grpc.NewServer(serverOpts...)
	server := &rpcServer{
		opts:               d.opts,
		logger:             d.logger.Named("rpc_server"),
		session:            d.mnemosyne,
		passwordHasher:     passwordHasher,
		permissionRegistry: permissionReg,
		repository:         repos,
	}

	charonrpc.RegisterAuthServer(gRPCServer, newAuth(server))
	charonrpc.RegisterUserManagerServer(gRPCServer, newUserManager(server))
	charonrpc.RegisterGroupManagerServer(gRPCServer, newGroupManager(server))
	charonrpc.RegisterPermissionManagerServer(gRPCServer, newPermissionManager(server))
	charonrpc.RegisterRefreshTokenManagerServer(gRPCServer, newRefreshTokenManager(server))

	if !d.opts.Test {
		prometheus.DefaultRegisterer.Register(interceptor)
		promgrpc.RegisterInterceptor(gRPCServer, interceptor)
	}

	go func() {
		d.logger.Info("rpc server is running", zap.Stringer("address", d.rpcListener.Addr()))

		if err := gRPCServer.Serve(d.rpcListener); err != nil {
			if err == grpc.ErrServerStopped {
				d.logger.Info("grpc server has been stopped")
				return
			}

			d.logger.Error("grpc server stopped with an error", zap.Error(err))
		}
	}()

	if d.debugListener != nil {
		go func() {
			d.logger.Info("debug server is running", zap.Stringer("address", d.debugListener.Addr()))
			// TODO: implement keep alive

			mux := http.NewServeMux()
			mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
			mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
			mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
			mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
			mux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
			mux.Handle("/healthz", &healthHandler{
				logger:   d.logger,
				postgres: db,
			})
			mux.Handle("/debug/requests", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				//any, sensitive := trace.AuthRequest(req)
				//if !any {
				//	http.Error(w, "not allowed", http.StatusUnauthorized)
				//	return
				//}
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				trace.Render(w, req, true)
			}))
			mux.Handle("/debug/events", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				//any, sensitive := trace.AuthRequest(req)
				//if !any {
				//	http.Error(w, "not allowed", http.StatusUnauthorized)
				//	return
				//}
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				trace.RenderEvents(w, req, true)
			}))
			if err := http.Serve(d.debugListener, mux); err != nil {
				d.logger.Error("debug server stopped with an error", zap.Error(err))
			}
		}()
	}

	return
}

// Close implements io.Closer interface.
func (d *Daemon) Close() (err error) {
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
