package main

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type getGroupHandler struct {
	*handler
}

func (ggh *getGroupHandler) handle(ctx context.Context, req *charon.GetGroupRequest) (*charon.GetGroupResponse, error) {
	ggh.loggerWith("group_id", req.Id)

	act, err := ggh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = ggh.firewall(req, act); err != nil {
		return nil, err
	}

	entity, err := ggh.repository.group.FindOneByID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "charond: group with id %d does not exists", req.Id)
		}
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	return &charon.GetGroupResponse{
		Group: entity.Message(),
	}, nil
}

func (ggh *getGroupHandler) firewall(req *charon.GetGroupRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.GroupCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charond: group cannot be retrieved, missing permission")
}
