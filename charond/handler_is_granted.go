package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type isGrantedHandler struct {
	*handler
}

func (ig *isGrantedHandler) handle(ctx context.Context, req *charon.IsGrantedRequest) (*charon.IsGrantedResponse, error) {
	ig.loggerWith("user_id", req.UserId, "permission", req.Permission)

	if req.Permission == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "permission cannot be empty")
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

	granted, err := ig.repository.userPermissions.IsGranted(req.UserId, charon.Permission(req.Permission))
	if err != nil {
		return nil, err
	}

	return &charon.IsGrantedResponse{
		Granted: granted,
	}, nil
}

func (ig *isGrantedHandler) firewall(req *charon.IsGrantedRequest, act *actor) error {
	if act.user.ID == req.UserId {
		return nil
	}
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.UserPermissionCanCheckGrantingAsStranger) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charond: group granting cannot be checked, missing permission")
}
