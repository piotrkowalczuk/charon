package main

import (
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type rpcServer struct {
	meta           metadata.MD
	logger         log.Logger
	monitor        *monitoring
	session        mnemosyne.Mnemosyne
	passwordHasher charon.PasswordHasher
	repository     struct {
		user       UserRepository
		permission PermissionRepository
		group      GroupRepository
	}
}

func (rs *rpcServer) retrieveActor(ctx context.Context, token mnemosyne.Token) (user *userEntity, session *mnemosyne.Session, permissions charon.Permissions, err error) {
	var userID int64
	var entities []*permissionEntity

	session, err = rs.session.Get(ctx, token)
	if err != nil {
		if err == mnemosyne.ErrSessionNotFound {
			err = grpc.Errorf(codes.Unauthenticated, "charond: action cannot be performed: %s", grpc.ErrorDesc(err))
			return
		}
		return
	}

	userID, err = charon.SessionSubjectID(session.SubjectId).UserID()
	if err != nil {
		return
	}
	user, err = rs.repository.user.FindOneByID(userID)
	if err != nil {
		return
	}

	entities, err = rs.repository.permission.FindByUserID(userID)

	permissions = make(charon.Permissions, 0, len(entities))
	for _, e := range entities {
		permissions = append(permissions, e.Permission())
	}

	return
}

func (rs *rpcServer) token(ctx context.Context) (mnemosyne.Token, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return mnemosyne.Token{}, errors.New("charond: missing metadata in context, session token cannot be retrieved")
	}

	if len(md[mnemosyne.TokenMetadataKey]) == 0 {
		return mnemosyne.Token{}, errors.New("charond: missing sesion token in metadata")
	}

	return mnemosyne.DecodeToken(md[mnemosyne.TokenMetadataKey][0]), nil
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
