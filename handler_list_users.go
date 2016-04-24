package charon

import (
	"github.com/piotrkowalczuk/protot"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listUsersHandler struct {
	*handler
}

func (luh *listUsersHandler) handle(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
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
	if act.permissions.Contains(UserCanRetrieveAsOwner, UserCanRetrieveStaffAsOwner) {
		criteria.createdBy = &protot.QueryInt64{Values: []int64{act.user.ID}, Type: protot.NumericQueryType_EQUAL, Valid: true}
	}

	ents, err := luh.repository.user.Find(criteria)
	if err != nil {
		return nil, err
	}
	luh.loggerWith("count", len(ents))

	return luh.response(ents)
}

func (luh *listUsersHandler) firewall(req *ListUsersRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if req.IsSuperuser.BoolOr(false) {
		return grpc.Errorf(codes.PermissionDenied, "charon: only superuser is permited to retrieve other superusers")
	}
	if req.CreatedBy == nil {
		if !act.permissions.Contains(UserCanRetrieveAsStranger) {
			return grpc.Errorf(codes.PermissionDenied, "charon: list of users cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	if req.IsStaff.BoolOr(false) {
		if req.CreatedBy.Value() == act.user.ID {
			if !act.permissions.Contains(UserCanRetrieveStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "charon: list of staff users cannot be retrieved as an owner, missing permission")
			}
			return nil
		}
		if !act.permissions.Contains(UserCanRetrieveStaffAsStranger) {
			return grpc.Errorf(codes.PermissionDenied, "charon: list of staff users cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	if req.CreatedBy.Value() == act.user.ID {
		if !act.permissions.Contains(UserCanRetrieveAsOwner) {
			return grpc.Errorf(codes.PermissionDenied, "charon: list of users cannot be retrieved as an owner, missing permission")
		}
		return nil
	}

	return nil
}

func (luh *listUsersHandler) response(ents []*userEntity) (*ListUsersResponse, error) {
	resp := &ListUsersResponse{
		Users: make([]*User, 0, len(ents)),
	}
	var (
		err error
		msg *User
	)
	for _, e := range ents {
		if msg, err = e.Message(); err != nil {
			return nil, err
		}
		resp.Users = append(resp.Users, msg)
	}

	return resp, nil
}
