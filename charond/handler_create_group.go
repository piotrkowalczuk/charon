package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type createGroupHandler struct {
	*handler
}

func (cgh *createGroupHandler) handle(ctx context.Context, req *charon.CreateGroupRequest) (*charon.CreateGroupResponse, error) {
	actor, err := cgh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	if !actor.permissions.Contains(charon.GroupCanCreate) {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: actor do not have permission: %s", charon.GroupCanCreate.String())
	}

	entity, err := cgh.repository.group.Create(actor.user.ID, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return &charon.CreateGroupResponse{
		Group: entity.Message(),
	}, nil
}
