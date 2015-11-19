package main

import (
	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
)

type rpcServer struct {
	logger               log.Logger
	monitor              *monitoring
	session              mnemosyne.Mnemosyne
	passwordHasher       charon.PasswordHasher
	userRepository       UserRepository
	permissionRepository PermissionRepository
}

func (rs *rpcServer) retrieveUserData(ctx context.Context) (user *userEntity, session *mnemosyne.Session, permissions charon.Permissions, err error) {
	var userID int64
	var entities []*permissionEntity

	session, err = rs.session.Get(ctx)
	if err != nil {
		return
	}

	userID, err = charon.UserIDFromSession(session)
	if err != nil {
		return
	}

	user, err = rs.userRepository.FindOneByID(userID)
	if err != nil {
		return
	}

	entities, err = rs.permissionRepository.FindByUserID(userID)

	permissions = make(charon.Permissions, 0, len(entities))
	for _, e := range entities {
		permissions = append(permissions, e.Permission())
	}

	return
}
