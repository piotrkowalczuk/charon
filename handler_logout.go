package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type logoutHandler struct {
	*handler
}

func (lh *logoutHandler) handle(ctx context.Context, r *LogoutRequest) (*LogoutResponse, error) {
	if r.AccessToken.IsEmpty() { // TODO: probably wrong, implement IsEmpty method for ID
		return nil, grpc.Errorf(codes.InvalidArgument, "charon: empty session id, logout aborted")
	}

	if err := lh.session.Abandon(ctx, *r.AccessToken); err != nil {
		return nil, err
	}

	lh.loggerWith("token", r.AccessToken.Encode())

	return &LogoutResponse{}, nil
}
