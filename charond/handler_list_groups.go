package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listGroupsHandler struct {
	*handler
}

func (lgh *listGroupsHandler) handle(ctx context.Context, req *charon.ListGroupsRequest) (*charon.ListGroupsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: list groups endpoint is not implemented yet")
}
