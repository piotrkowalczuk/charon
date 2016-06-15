package charond

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
	if r.AccessToken.IsEmpty() { // TODO: probably wrong, implement IsEmpty method for ID
		return nil, grpc.Errorf(codes.InvalidArgument, "empty session id, logout aborted")
	}

	if err := lh.session.Abandon(ctx, r.AccessToken.Encode()); err != nil {
		return nil, err
	}

	lh.loggerWith("token", r.AccessToken.Encode())

	return &charon.LogoutResponse{}, nil
}
