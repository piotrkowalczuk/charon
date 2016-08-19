package charond

import (
	"strings"

	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/go-ldap/ldap"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
)

type handler struct {
	opts       DaemonOpts
	logger     log.Logger
	monitor    monitoringRPC
	session    mnemosynerpc.SessionManagerClient
	repository repositories
	ldap       *ldap.Conn
}

func newHandler(rs *rpcServer, ctx context.Context, endpoint string) *handler {
	h := &handler{
		opts:       rs.opts,
		logger:     rs.loggerBackground(ctx, "endpoint", endpoint),
		session:    rs.session,
		repository: rs.repository,
		ldap:       rs.ldap,
	}

	return h
}

func handleMnemosyneError(err error) error {
	if grpc.Code(err) == codes.NotFound {
		return grpc.Errorf(codes.Unauthenticated, grpc.ErrorDesc(err))
	}

	return err
}

func (h *handler) loggerWith(keyval ...interface{}) {
	h.logger = log.NewContext(h.logger).With(keyval...)
}

func (h *handler) retrieveActor(ctx context.Context) (act *actor, err error) {
	var (
		userID   int64
		entities []*permissionEntity
		res      *mnemosynerpc.ContextResponse
	)

	res, err = h.session.Context(ctx, none())
	if err != nil {
		// TODO: make it better ;(
		if peer, ok := peer.FromContext(ctx); ok {
			if strings.HasPrefix(peer.Addr.String(), "127.0.0.1") {
				return &actor{
					user:    &userEntity{},
					isLocal: true,
				}, nil
			}
		}
		err = handleMnemosyneError(err)
		return
	}
	userID, err = charon.SubjectID(res.Session.SubjectId).UserID()
	if err != nil {
		return
	}

	act = &actor{}
	act.user, err = h.repository.user.findOneByID(userID)
	if err != nil {
		return
	}
	entities, err = h.repository.permission.findByUserID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return act, nil
		}
		return
	}

	act.permissions = make(charon.Permissions, 0, len(entities))
	for _, e := range entities {
		act.permissions = append(act.permissions, e.Permission())
	}

	return
}
