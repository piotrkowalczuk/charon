package main

import (
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type rpcServer struct {
	meta               metadata.MD
	logger             log.Logger
	monitor            *monitoring
	session            mnemosyne.Mnemosyne
	passwordHasher     charon.PasswordHasher
	permissionRegistry PermissionRegistry
	repository         struct {
		user       UserRepository
		permission PermissionRepository
		group      GroupRepository
	}
}

type actor struct {
	user        *userEntity
	session     *mnemosyne.Session
	permissions charon.Permissions
}

func (a *actor) isLocalhost() bool {
	return a.user == nil && a.session == nil && a.permissions == nil
}

func (rs *rpcServer) retrieveActor(ctx context.Context) (a *actor, err error) {
	var (
		userID   int64
		entities []*permissionEntity
		ses      *mnemosyne.Session
	)

	ses, err = rs.session.FromContext(ctx)
	if err != nil {
		if peer, ok := peer.FromContext(ctx); ok {
			if strings.HasPrefix(peer.Addr.String(), "127.0.0.1") {
				return &actor{}, nil
			}
		}
		return
	}
	userID, err = charon.SessionSubjectID(ses.SubjectId).UserID()
	if err != nil {
		return
	}
	a.user, err = rs.repository.user.FindOneByID(userID)
	if err != nil {
		return
	}

	entities, err = rs.repository.permission.FindByUserID(userID)
	if err != nil {
		return
	}

	a.permissions = make(charon.Permissions, 0, len(entities))
	for _, e := range entities {
		a.permissions = append(a.permissions, e.Permission())
	}

	return
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
