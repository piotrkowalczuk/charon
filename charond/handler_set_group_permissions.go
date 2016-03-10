package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type setGroupPermissionsHandler struct {
	*handler
}

func (sgh *setGroupPermissionsHandler) handle(ctx context.Context, req *charon.SetGroupPermissionsRequest) (*charon.SetGroupPermissionsResponse, error) {
	sgh.loggerWith("group_id", req.GroupId)
	return nil, grpc.Errorf(codes.Unimplemented, "charond: set group permissions endpoint is not implemented yet")
}
