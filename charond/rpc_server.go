package charond

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type rpcServer struct {
	opts               DaemonOpts
	meta               metadata.MD
	logger             log.Logger
	ldap               *sync.Pool
	session            mnemosynerpc.SessionManagerClient
	passwordHasher     password.Hasher
	permissionRegistry model.PermissionRegistry
	repository         repositories
}

func (rs *rpcServer) loggerBackground(ctx context.Context, keyval ...interface{}) log.Logger {
	l := log.NewContext(rs.logger).With(keyval...)
	if md, ok := metadata.FromContext(ctx); ok {
		if rid, ok := md["request_id"]; ok && len(rid) >= 1 {
			l = l.With("request_id", rid[0])
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		l = l.With("peer_address", p.Addr.String())
	}

	return l
}

// Context create new context based on given metadata and instance metadata.
func (rs *rpcServer) Context(md metadata.MD) context.Context {
	if md.Len() == 0 {
		md = rs.meta
	} else {
		md = rs.metadata(md)
	}

	return metadata.NewContext(context.Background(), md)
}

func (rs *rpcServer) metadata(md metadata.MD) metadata.MD {
	for key, value := range rs.meta {
		if _, ok := md[key]; !ok {
			md[key] = value
		}
	}

	return md
}
