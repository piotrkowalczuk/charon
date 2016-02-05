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
