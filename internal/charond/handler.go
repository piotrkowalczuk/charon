package charond

import (
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
)

type handler struct {
	session.ActorProvider

	opts       DaemonOpts
	logger     *zap.Logger
	repository repositories
	session    mnemosynerpc.SessionManagerClient
}

func newHandler(rs *rpcServer) *handler {
	h := &handler{
		opts:       rs.opts,
		session:    rs.session,
		repository: rs.repository,
		logger:     rs.logger,
		ActorProvider: &session.MnemosyneActorProvider{
			Client:             rs.session,
			UserProvider:       rs.repository.user,
			PermissionProvider: rs.repository.permission,
		},
	}

	return h
}
