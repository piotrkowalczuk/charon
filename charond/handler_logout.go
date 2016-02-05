package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type logoutHandler struct {
	*handler
}

func (lh *logoutHandler) handle(ctx context.Context, r *charon.LogoutRequest) (*charon.LogoutResponse, error) {
	if r.Token.IsEmpty() { // TODO: probably wrong, implement IsEmpty method for ID
		return nil, grpc.Errorf(codes.InvalidArgument, "charond: empty session id, logout aborted")
	}

	if err := lh.session.Abandon(ctx, *r.Token); err != nil {
		return nil, err
	}

	lh.loggerWith("token", r.Token.Encode())

	return &charon.LogoutResponse{}, nil
}
