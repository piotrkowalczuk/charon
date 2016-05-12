package charond

import (
	"errors"
	"strings"

	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
)

type handler struct {
	logger     log.Logger
	monitor    monitoringRPC
	session    mnemosyne.Mnemosyne
	repository repositories
}

func newHandler(rs *rpcServer, ctx context.Context, endpoint string) *handler {
	h := &handler{
		logger:     rs.loggerBackground(ctx, "endpoint", endpoint),
		session:    rs.session,
		repository: rs.repository,
	}
	if rs.monitor.enabled {
		h.monitor = rs.monitor.rpc.with(metrics.Field{Key: "endpoint", Value: endpoint})
	}

	return h
}

func (h *handler) handle(err error, msg string) error {
	if err != nil {
		if h.monitor.enabled {
			h.monitor.errors.With(metrics.Field{Key: "code", Value: grpc.Code(err).String()}).Add(1)
		}
		sklog.Error(h.logger, errors.New(grpc.ErrorDesc(err)))

		return grpc.Errorf(grpc.Code(err), "charond: %s", grpc.ErrorDesc(err))
	}

	sklog.Debug(h.logger, msg)
	return nil
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
		ses      *mnemosyne.Session
	)

	ses, err = h.session.FromContext(ctx)
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
	userID, err = charon.SubjectID(ses.SubjectId).UserID()
	if err != nil {
		return
	}

	act = &actor{}
	act.user, err = h.repository.user.FindOneByID(userID)
	if err != nil {
		return
	}
	entities, err = h.repository.permission.FindByUserID(userID)
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

func (h *handler) addRequest(i uint64) {
	if h.monitor.enabled {
		h.monitor.requests.Add(i)
	}
}
