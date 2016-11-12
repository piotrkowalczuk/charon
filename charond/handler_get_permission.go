package charond

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type getPermissionHandler struct {
	*handler
}

func (gph *getPermissionHandler) Get(ctx context.Context, req *charonrpc.GetPermissionRequest) (*charonrpc.GetPermissionResponse, error) {
	gph.loggerWith("permission_id", req.Id)

	if req.Id < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "permission id needs to be greater than zero")
	}

	act, err := gph.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = gph.firewall(req, act); err != nil {
		return nil, err
	}

	permission, err := gph.repository.permission.FindOneByID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "permission does not exists")
		}
		return nil, err
	}

	return &charonrpc.GetPermissionResponse{
		Permission: permission.Permission().String(),
	}, nil
}

func (gph *getPermissionHandler) firewall(req *charonrpc.GetPermissionRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.PermissionCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "permission cannot be retrieved, missing permission")
}
