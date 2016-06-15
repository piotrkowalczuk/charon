package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/qtypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listUsersHandler struct {
	*handler
}

func (luh *listUsersHandler) handle(ctx context.Context, req *charon.ListUsersRequest) (*charon.ListUsersResponse, error) {
	act, err := luh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = luh.firewall(req, act); err != nil {
		return nil, err
	}

	criteria := &userCriteria{
		sort:        req.Sort,
		offset:      req.Offset.Int64Or(0),
		limit:       req.Limit.Int64Or(10),
		isSuperuser: req.IsSuperuser,
		isStaff:     req.IsStaff,
		createdBy:   req.CreatedBy,
	}
	if act.permissions.Contains(charon.UserCanRetrieveAsOwner, charon.UserCanRetrieveStaffAsOwner) {
		criteria.createdBy = &qtypes.Int64{Values: []int64{act.user.ID}, Type: qtypes.NumericQueryType_EQUAL, Valid: true}
	}

	ents, err := luh.repository.user.Find(criteria)
	if err != nil {
		return nil, err
	}
	luh.loggerWith("count", len(ents))

	return luh.response(ents)
}

func (luh *listUsersHandler) firewall(req *charon.ListUsersRequest, act *actor) error {
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

func (luh *listUsersHandler) response(ents []*userEntity) (*charon.ListUsersResponse, error) {
	resp := &charon.ListUsersResponse{
		Users: make([]*charon.User, 0, len(ents)),
	}
	var (
		err error
		msg *charon.User
	)
	for _, e := range ents {
		if msg, err = e.message(); err != nil {
			return nil, err
		}
		resp.Users = append(resp.Users, msg)
	}

	return resp, nil
}
