package charond

import (
	"context"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type isAuthenticatedHandler struct {
	*handler
}

func (iah *isAuthenticatedHandler) IsAuthenticated(ctx context.Context, req *charonrpc.IsAuthenticatedRequest) (*wrappers.BoolValue, error) {
	if req.AccessToken == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "authentication status cannot be checked, missing token")
	}

	ses, err := iah.session.Get(ctx, &mnemosynerpc.GetRequest{AccessToken: req.AccessToken})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return &wrappers.BoolValue{Value: false}, nil
		}
		return nil, err
	}
	uid, err := session.ActorID(ses.Session.SubjectId).UserID()
	if err != nil {
		return nil, err
	}
	exists, err := iah.repository.user.Exists(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &wrappers.BoolValue{Value: exists}, nil
}
