package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type isAuthenticatedHandler struct {
	*handler
}

func (iah *isAuthenticatedHandler) handle(ctx context.Context, req *charon.IsAuthenticatedRequest) (*charon.IsAuthenticatedResponse, error) {
	if req.Token == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "charond: authentication status cannot be checked, missing token")
	}

	iah.loggerWith("token", req.Token.Encode())

	ok, err := iah.session.Exists(ctx, *req.Token)
	if err != nil {
		return nil, err
	}

	return &charon.IsAuthenticatedResponse{
		Authenticated: ok,
	}, nil
}
