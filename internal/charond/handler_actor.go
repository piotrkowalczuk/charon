package charond

import (
	"context"
	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"google.golang.org/grpc/codes"
)

type actorHandler struct {
	*handler
}

func (sh *actorHandler) Actor(ctx context.Context, r *wrappers.StringValue) (*charonrpc.ActorResponse, error) {
	var ses *mnemosynerpc.Session

	if r.Value == "" {
		res, err := sh.session.Context(ctx, none())
		if err != nil {
			return nil, handleMnemosyneError(err)
		}
		ses = res.Session
	} else {
		res, err := sh.session.Get(ctx, &mnemosynerpc.GetRequest{
			AccessToken: r.Value,
		})
		if err != nil {
			return nil, handleMnemosyneError(err)
		}
		ses = res.Session
	}

	id, err := session.ActorID(ses.SubjectId).UserID()
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "invalid session actor id")
	}

	ent, err := sh.repository.user.FindOneByID(ctx, id)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, grpcerr.E(codes.NotFound, "actor does not exists for given id")
	default:
		return nil, grpcerr.E(codes.Internal, "actor retrieval failure", err)
	}

	permissionEntities, err := sh.repository.permission.FindByUserID(ctx, id)
	switch err {
	case nil, sql.ErrNoRows:
	default:
		return nil, grpcerr.E(codes.Internal, "actor list of permissions failure", err)
	}

	permissions := make([]string, 0, len(permissionEntities))
	for _, e := range permissionEntities {
		permissions = append(permissions, e.Permission().String())
	}

	return &charonrpc.ActorResponse{
		Id:          int64(ent.ID),
		Username:    ent.Username,
		FirstName:   ent.FirstName,
		LastName:    ent.LastName,
		Permissions: permissions,
		IsActive:    ent.IsActive,
		IsConfirmed: ent.IsConfirmed,
		IsStuff:     ent.IsStaff,
		IsStaff:     ent.IsStaff,
		IsSuperuser: ent.IsSuperuser,
	}, nil
}
