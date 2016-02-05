package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type setUserGroupsHandler struct {
	*handler
}

func (sugh *setUserGroupsHandler) handle(ctx context.Context, req *charon.SetUserGroupsRequest) (*charon.SetUserGroupsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: set user groups endpoint is not implemented yet")
}
