package charond

import (
	"database/sql"
	"fmt"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type actorHandler struct {
	*handler
}

func (sh *actorHandler) Actor(ctx context.Context, r *wrappers.StringValue) (*charonrpc.ActorResponse, error) {
	var (
		ses *mnemosynerpc.Session
		err error
	)
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
		return nil, fmt.Errorf("invalid Session subject id: %s", ses.SubjectId)
	}

	ent, err := sh.repository.user.FindOneByID(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, grpc.Errorf(codes.NotFound, "subject does not exists with id: %d", id)
		case context.DeadlineExceeded, context.Canceled:
			return nil, err
		default:
			return nil, grpc.Errorf(codes.Internal, "subject retrieval failure: %s", err.Error())
		}
	}

	permissionEntities, err := sh.repository.permission.FindByUserID(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		switch err {
		case sql.ErrNoRows:
		case context.DeadlineExceeded, context.Canceled:
			return nil, err
		default:
			return nil, grpc.Errorf(codes.Internal, "subject list of permissions failure: %s", err.Error())
		}
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
		IsSuperuser: ent.IsSuperuser,
	}, nil
}
