package charond

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"

	"google.golang.org/grpc/codes"
)

type logoutHandler struct {
	*handler
}

func (lh *logoutHandler) Logout(ctx context.Context, r *charonrpc.LogoutRequest) (*empty.Empty, error) {
	if len(r.AccessToken) == 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "empty session id, logout aborted")
	}

	_, err := lh.session.Abandon(ctx, &mnemosynerpc.AbandonRequest{
		AccessToken: r.AccessToken,
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
