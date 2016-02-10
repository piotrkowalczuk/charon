package main

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/nilt"
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
		offset:      req.Offset.Int64Or(0),
		limit:       req.Limit.Int64Or(10),
		isSuperuser: nilBool(req.IsSuperuser),
		isStaff:     nilBool(req.IsStaff),
		createdBy:   nilInt64(req.CreatedBy),
	}
	if act.permissions.Contains(charon.UserCanRetrieveAsOwner, charon.UserCanRetrieveAsStranger) {
		criteria.createdBy = nilt.Int64{Int64: act.user.ID, Valid: true}
	}

	users, err := luh.repository.user.Find(criteria)
	if err != nil {
		return nil, err
	}

	luh.loggerWith("results", len(users))

	resp := &charon.ListUsersResponse{
		Users: make([]*charon.User, 0, len(users)),
	}
	for _, u := range users {
		resp.Users = append(resp.Users, u.Message())
	}

	return resp, nil
}

func (luh *listUsersHandler) firewall(req *charon.ListUsersRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if req.IsSuperuser.BoolOr(false) {
		return grpc.Errorf(codes.PermissionDenied, "charond: only superuser is permited to retrieve other superusers")
	}
	if req.IsStaff.BoolOr(false) {
		if req.CreatedBy.Int64Or(0) == act.user.ID {
			if !act.permissions.Contains(charon.UserCanRetrieveStaffAsOwner) {
				return grpc.Errorf(codes.PermissionDenied, "charond: list of staff users cannot be retrieved as an owner, missing permission")
			}
			return nil
		}
		if !act.permissions.Contains(charon.UserCanRetrieveStaffAsStranger) {
			return grpc.Errorf(codes.PermissionDenied, "charond: list of staff users cannot be retrieved as a stranger, missing permission")
		}
		return nil
	}
	if req.CreatedBy.Int64Or(0) == act.user.ID {
		if !act.permissions.Contains(charon.UserCanRetrieveAsOwner) {
			return grpc.Errorf(codes.PermissionDenied, "charond: list of users cannot be retrieved as an owner, missing permission")
		}
		return nil
	}
	if !act.permissions.Contains(charon.UserCanRetrieveAsStranger) {
		return grpc.Errorf(codes.PermissionDenied, "charond: list of users cannot be retrieved as a stranger, missing permission")
	}

	return nil
}
