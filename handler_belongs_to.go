package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type belongsToHandler struct {
	*handler
}

func (ig *belongsToHandler) handle(ctx context.Context, req *BelongsToRequest) (*BelongsToResponse, error) {
	ig.loggerWith("user_id", req.UserId, "group_id", req.GroupId)

	if req.GroupId < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "group id needs to be greater than zero")
	}
	if req.UserId < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "user id needs to be greater than zero")
	}

	act, err := ig.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = ig.firewall(req, act); err != nil {
		return nil, err
	}

	belongs, err := ig.repository.userGroups.Exists(req.UserId, req.GroupId)
	if err != nil {
		return nil, err
	}

	return &BelongsToResponse{
		Belongs: belongs,
	}, nil
}

func (ig *belongsToHandler) firewall(req *BelongsToRequest, act *actor) error {
	if act.user.ID == req.UserId {
		return nil
	}
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(UserGroupCanCheckBelongingAsStranger) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charon: group belonging cannot be checked, missing permission")
}
