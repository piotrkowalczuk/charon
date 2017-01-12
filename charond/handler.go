package charond

import (
	"strings"

	"database/sql"

	"sync"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
)

type handler struct {
	opts       DaemonOpts
	logger     log.Logger
	repository repositories
	session    mnemosynerpc.SessionManagerClient
	monitor    monitoringRPC
	ldap       *sync.Pool
}

func newHandler(rs *rpcServer) *handler {
	h := &handler{
		opts:       rs.opts,
		session:    rs.session,
		repository: rs.repository,
		ldap:       rs.ldap,
		logger:     rs.logger,
	}

	return h
}

func handleMnemosyneError(err error) error {
	if grpc.Code(err) == codes.NotFound {
		return grpc.Errorf(codes.Unauthenticated, grpc.ErrorDesc(err))
	}

	return err
}

func (h *handler) retrieveActor(ctx context.Context) (act *session.Actor, err error) {
	var (
		userID   int64
		entities []*model.PermissionEntity
		res      *mnemosynerpc.ContextResponse
	)

	res, err = h.session.Context(ctx, none())
	if err != nil {
		// TODO: make it better ;(
		if peer, ok := peer.FromContext(ctx); ok {
			if strings.HasPrefix(peer.Addr.String(), "127.0.0.1") {
				return &session.Actor{
					User:    &model.UserEntity{},
					IsLocal: true,
				}, nil
			}
		}
		err = handleMnemosyneError(err)
		return
	}
	userID, err = session.ActorID(res.Session.SubjectId).UserID()
	if err != nil {
		return
	}

	act = &session.Actor{}
	act.User, err = h.repository.user.FindOneByID(ctx, userID)
	if err != nil {
		return
	}
	entities, err = h.repository.permission.FindByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return act, nil
		}
		return
	}

	act.Permissions = make(charon.Permissions, 0, len(entities))
	for _, e := range entities {
		act.Permissions = append(act.Permissions, e.Permission())
	}

	return
}
