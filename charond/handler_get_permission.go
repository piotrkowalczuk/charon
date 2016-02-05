package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type getPermissionHandler struct {
	*handler
}

func (gph *getPermissionHandler) handle(ctx context.Context, req *charon.GetPermissionRequest) (*charon.GetPermissionResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "charond: get permission endpoint is not implemented yet")
}
