package main

import (
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type handler struct {
	logger     log.Logger
	monitor    monitoringRPC
	session    mnemosyne.Mnemosyne
	repository struct {
		user       UserRepository
		permission PermissionRepository
		group      GroupRepository
	}
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

func (h *handler) handle(err error, msg string) {
	if err != nil {
		code := grpc.Code(err)

		h.loggerWith("code", code)
		if h.monitor.enabled {
			h.monitor.errors.With(metrics.Field{Key: "code", Value: code.String()}).Add(1)
		}
		sklog.Error(h.logger, err)

		return
	}

	sklog.Debug(h.logger, msg)
}

func (h *handler) loggerWith(keyval ...interface{}) {
	h.logger = log.NewContext(h.logger).With(keyval...)
}

func (h *handler) retrieveActor(ctx context.Context) (a *actor, err error) {
	var (
		userID   int64
		entities []*permissionEntity
		ses      *mnemosyne.Session
	)

	ses, err = h.session.FromContext(ctx)
	if err != nil {
		if peer, ok := peer.FromContext(ctx); ok {
			if strings.HasPrefix(peer.Addr.String(), "127.0.0.1") {
				return &actor{}, nil
			}
		}
		return
	}
	userID, err = charon.SubjectID(ses.SubjectId).UserID()
	if err != nil {
		return
	}
	a.user, err = h.repository.user.FindOneByID(userID)
	if err != nil {
		return
	}

	entities, err = h.repository.permission.FindByUserID(userID)
	if err != nil {
		return
	}

	a.permissions = make(charon.Permissions, 0, len(entities))
	for _, e := range entities {
		a.permissions = append(a.permissions, e.Permission())
	}

	return
}

func (h *handler) addRequest(i uint64) {
	if h.monitor.enabled {
		h.monitor.requests.Add(i)
	}
}
