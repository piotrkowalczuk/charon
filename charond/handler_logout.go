package charond

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type logoutHandler struct {
	*handler
}

func (lh *logoutHandler) Logout(ctx context.Context, r *charonrpc.LogoutRequest) (*empty.Empty, error) {
	if len(r.AccessToken) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "empty session id, logout aborted")
	}

	_, err := lh.session.Abandon(ctx, &mnemosynerpc.AbandonRequest{
		AccessToken: r.AccessToken,
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
