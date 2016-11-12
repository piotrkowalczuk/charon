package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
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

	criteria := &model.UserCriteria{
		Sort:        req.Sort,
		Offset:      req.Offset.Int64Or(0),
		Limit:       req.Limit.Int64Or(10),
		IsSuperuser: req.IsSuperuser,
		IsStaff:     req.IsStaff,
		CreatedBy:   req.CreatedBy,
	}
	if act.permissions.Contains(charon.UserCanRetrieveAsOwner, charon.UserCanRetrieveStaffAsOwner) {
		criteria.CreatedBy = qtypes.EqualInt64(act.user.ID)
	}

	ents, err := luh.repository.user.Find(criteria)
	if err != nil {
		return nil, err
	}
	luh.loggerWith("count", len(ents))

	return luh.response(ents)
}

func (luh *listUsersHandler) firewall(req *charonrpc.ListUsersRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if req.IsSuperuser.BoolOr(false) {
		return grpc.Errorf(codes.PermissionDenied, "only superuser is permited to retrieve other superusers")
	}
	if req.CreatedBy == nil {
		if !act.permissions.Contains(charon.UserCanRetrieveAsStranger) {
			return grpc.Errorf(codes.PermissionDenied, "list of users cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	if req.IsStaff.BoolOr(false) {
		if req.CreatedBy.Value() == act.user.ID {
			if !act.permissions.Contains(charon.UserCanRetrieveStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "list of staff users cannot be retrieved as an owner, missing permission")
			}
			return nil
		}
		if !act.permissions.Contains(charon.UserCanRetrieveStaffAsStranger) {
			return grpc.Errorf(codes.PermissionDenied, "list of staff users cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	if req.CreatedBy.Value() == act.user.ID {
		if !act.permissions.Contains(charon.UserCanRetrieveAsOwner) {
			return grpc.Errorf(codes.PermissionDenied, "list of users cannot be retrieved as an owner, missing permission")
		}
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
