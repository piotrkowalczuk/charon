package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type logoutHandler struct {
	*handler
}

func (lh *logoutHandler) handle(ctx context.Context, r *charon.LogoutRequest) (*charon.LogoutResponse, error) {
	if len(r.AccessToken) == 0 { // TODO: probably wrong, implement IsEmpty method for ID
		return nil, grpc.Errorf(codes.InvalidArgument, "empty session id, logout aborted")
	}

	_, err := lh.session.Abandon(ctx, &mnemosynerpc.AbandonRequest{
		AccessToken: r.AccessToken,
	})
	if err != nil {
		return nil, err
	}

	lh.loggerWith("token", r.AccessToken)

	return &charon.LogoutResponse{}, nil
}
