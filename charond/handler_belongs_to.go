package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type belongsToHandler struct {
	*handler
}

func (ig *belongsToHandler) handle(ctx context.Context, r *charon.BelongsToRequest) (*charon.BelongsToResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "belongs to endpoint is not implemented yet")
}
