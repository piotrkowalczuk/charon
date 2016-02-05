package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listGroupPermissionsHandler struct {
	*handler
}

func (lgh *listGroupPermissionsHandler) handle(ctx context.Context, req *charon.ListGroupPermissionsRequest) (*charon.ListGroupPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: list group permissions endpoint is not implemented yet")
}
