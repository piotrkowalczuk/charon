package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/qtypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listUsersHandler struct {
	*handler
}

func (luh *listUsersHandler) List(ctx context.Context, req *charonrpc.ListUsersRequest) (*charonrpc.ListUsersResponse, error) {
	act, err := luh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = luh.firewall(req, act); err != nil {
		return nil, err
	}

	cri := &model.UserCriteria{
		Sort:        req.Sort,
		Offset:      req.Offset.Int64Or(0),
		Limit:       req.Limit.Int64Or(10),
		IsSuperuser: *req.IsSuperuser,
		IsStaff:     *req.IsStaff,
		CreatedBy:   req.CreatedBy,
	}

	if !act.User.IsSuperuser {
		cri.IsSuperuser = ntypes.False()
	}
	if !act.Permissions.Contains(charon.UserCanRetrieveStaffAsStranger) {
		cri.IsStaff = ntypes.False()
	}
	if act.Permissions.Contains(charon.UserCanRetrieveAsOwner, charon.UserCanRetrieveStaffAsOwner) {
		cri.CreatedBy = qtypes.EqualInt64(act.User.ID)
	}

	ents, err := luh.repository.user.Find(ctx, cri)
	if err != nil {
		return nil, err
	}
	return luh.response(ents)
}

func (luh *listUsersHandler) firewall(req *charonrpc.ListUsersRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if req.IsSuperuser.BoolOr(false) {
		return grpc.Errorf(codes.PermissionDenied, "only superuser is permited to retrieve other superusers")
	}
	if req.CreatedBy == nil {
		if !act.Permissions.Contains(charon.UserCanRetrieveAsStranger) {
			return grpc.Errorf(codes.PermissionDenied, "list of users cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	if req.IsStaff.BoolOr(false) {
		if req.CreatedBy.Value() == act.User.ID {
			if !act.Permissions.Contains(charon.UserCanRetrieveStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "list of staff users cannot be retrieved as an owner, missing permission")
			}
			return nil
		}
		if !act.Permissions.Contains(charon.UserCanRetrieveStaffAsStranger) {
			return grpc.Errorf(codes.PermissionDenied, "list of staff users cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	if req.CreatedBy.Value() == act.User.ID {
		if !act.Permissions.Contains(charon.UserCanRetrieveAsStranger, charon.UserCanRetrieveAsOwner) {
			return grpc.Errorf(codes.PermissionDenied, "list of users cannot be retrieved as an owner, missing permission")
		}
		return nil
	}
	if !act.Permissions.Contains(charon.UserCanRetrieveAsStranger) {
		return nil
	}

	return nil
}

func (luh *listUsersHandler) response(ents []*model.UserEntity) (*charonrpc.ListUsersResponse, error) {
	resp := &charonrpc.ListUsersResponse{
		Users: make([]*charonrpc.User, 0, len(ents)),
	}
	var (
		err error
		msg *charonrpc.User
	)
	for _, e := range ents {
		if msg, err = e.Message(); err != nil {
			return nil, err
		}
		resp.Users = append(resp.Users, msg)
	}

	return resp, nil
}
