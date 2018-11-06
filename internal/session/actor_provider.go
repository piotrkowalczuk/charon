package session

import (
	"context"
	"database/sql"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type ActorProvider interface {
	Actor(context.Context) (*Actor, error)
}

type MnemosyneActorProvider struct {
	Client             mnemosynerpc.SessionManagerClient
	UserProvider       model.UserProvider
	PermissionProvider model.PermissionProvider
}

func (p *MnemosyneActorProvider) Actor(ctx context.Context) (*Actor, error) {
	var (
		act      *Actor
		userID   int64
		entities []*model.PermissionEntity
		res      *mnemosynerpc.ContextResponse
	)

	res, err := p.Client.Context(ctx, &empty.Empty{})
	if err != nil {
		if isLocal(ctx) {
			return &Actor{
				User:    &model.UserEntity{},
				IsLocal: true,
			}, nil
		}
		return nil, handleMnemosyneError(err)
	}

	userID, err = ActorID(res.Session.SubjectId).UserID()
	if err != nil {
		return nil, grpcerr.E(codes.InvalidArgument, err)
	}

	act = &Actor{}
	act.User, err = p.UserProvider.FindOneByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpcerr.E(codes.PermissionDenied, "actor does not exists")
		}
		return nil, grpcerr.E(codes.Internal, "actor fetch failure", err)
	}
	entities, err = p.PermissionProvider.FindByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return act, nil
		}
		return nil, grpcerr.E(codes.Internal, "permissions fetch failure", err)
	}

	act.Permissions = make(charon.Permissions, 0, len(entities))
	for _, e := range entities {
		act.Permissions = append(act.Permissions, e.Permission())
	}

	return act, nil
}

func isLocal(ctx context.Context) bool {
	if p, ok := peer.FromContext(ctx); ok {
		if strings.HasPrefix(p.Addr.String(), "127.0.0.1") {
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				if len(md["user-agent"]) == 1 && strings.HasPrefix(md["user-agent"][0], "charonctl") {
					return true
				}
			}
		}
	}
	return false
}

func handleMnemosyneError(err error) error {
	if sts, ok := status.FromError(err); ok {
		switch sts.Code() {
		case codes.NotFound:
			return grpcerr.E(codes.Unauthenticated, "session not found")
		case codes.InvalidArgument:
			return grpcerr.E(codes.Unauthenticated, sts.Message())
		}
	}

	return grpcerr.E(codes.Internal, "session fetch failure", err)
}
