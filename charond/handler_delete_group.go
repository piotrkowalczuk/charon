package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
)

type deleteGroupHandler struct {
	*handler
}

func (dgh *deleteGroupHandler) handle(ctx context.Context, req *charon.DeleteGroupRequest) (*charon.DeleteGroupResponse, error) {
	dgh.loggerWith("group_id", req.Id)

	affected, err := dgh.repository.group.DeleteOneByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &charon.DeleteGroupResponse{
		Affected: affected,
	}, nil
}
