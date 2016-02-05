package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type setUserPermissionsHandler struct {
	*handler
}

func (suph *setUserPermissionsHandler) handle(ctx context.Context, req *charon.SetUserPermissionsRequest) (*charon.SetUserPermissionsResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: set user permissions endpoint is not implemented yet")
}
