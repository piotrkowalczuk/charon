package charond

import (
	"context"

	"github.com/piotrkowalczuk/charon"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/qtypes"

	"google.golang.org/grpc/codes"
)

type listUsersHandler struct {
	*handler
}

func (luh *listUsersHandler) List(ctx context.Context, req *charonrpc.ListUsersRequest) (*charonrpc.ListUsersResponse, error) {
	act, err := luh.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = luh.firewall(req, act); err != nil {
		return nil, err
	}

	cri := &model.UserCriteria{
		IsSuperuser: allocNilBool(req.IsSuperuser),
		IsStaff:     allocNilBool(req.IsStaff),
		CreatedBy:   req.CreatedBy,
	}

	if !act.User.IsSuperuser {
		cri.IsSuperuser = *ntypes.False()
	}
	if !act.Permissions.Contains(charon.UserCanRetrieveStaffAsStranger) {
		cri.IsStaff = *ntypes.False()
	}
	if act.Permissions.Contains(charon.UserCanRetrieveAsOwner, charon.UserCanRetrieveStaffAsOwner) {
		cri.CreatedBy = qtypes.EqualInt64(act.User.ID)
	}

	ents, err := luh.repository.user.Find(ctx, &model.UserFindExpr{
		OrderBy: mapping.OrderBy(req.OrderBy),
		Offset:  req.Offset.Int64Or(0),
		Limit:   req.Limit.Int64Or(10),
		Where:   cri,
	})
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "find users query failed", err)
	}
	return luh.response(ents)
}

func (luh *listUsersHandler) firewall(req *charonrpc.ListUsersRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if req.IsSuperuser.BoolOr(false) {
		return grpcerr.E(codes.PermissionDenied, "only superuser is permitted to retrieve other superusers")
	}
	// STAFF USERS
	if req.IsStaff.BoolOr(false) {
		if req.CreatedBy != nil && req.CreatedBy.Value() == act.User.ID {
			if !act.Permissions.Contains(charon.UserCanRetrieveStaffAsStranger, charon.UserCanRetrieveStaffAsOwner) {
				return grpcerr.E(codes.PermissionDenied, "list of staff users cannot be retrieved as an owner, missing permission")
			}
			return nil
		}
		if !act.Permissions.Contains(charon.UserCanRetrieveStaffAsStranger) {
			return grpcerr.E(codes.PermissionDenied, "list of staff users cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	// NON STAFF USERS
	if req.CreatedBy != nil && req.CreatedBy.Value() == act.User.ID {
		if !act.Permissions.Contains(charon.UserCanRetrieveAsStranger, charon.UserCanRetrieveAsOwner) {
			return grpcerr.E(codes.PermissionDenied, "list of users cannot be retrieved as an owner, missing permission")
		}
		return nil
	}
	if !act.Permissions.Contains(charon.UserCanRetrieveAsStranger) {
		return grpcerr.E(codes.PermissionDenied, "list of users cannot be retrieved as a stranger, missing permission")
	}
	return nil
}

func (luh *listUsersHandler) response(ents []*model.UserEntity) (*charonrpc.ListUsersResponse, error) {
	msg, err := mapping.ReverseUsers(ents)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "user reverse mapping failure")
	}
	return &charonrpc.ListUsersResponse{
		Users: msg,
	}, nil
}
