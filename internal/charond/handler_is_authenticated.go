package charond

import (
	"context"
	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type isAuthenticatedHandler struct {
	*handler
}

func (iah *isAuthenticatedHandler) IsAuthenticated(ctx context.Context, req *charonrpc.IsAuthenticatedRequest) (*wrappers.BoolValue, error) {
	if req.AccessToken == "" {
		return nil, grpcerr.E(codes.InvalidArgument, "authentication status cannot be checked, missing access token")
	}

	ses, err := iah.session.Get(ctx, &mnemosynerpc.GetRequest{AccessToken: req.AccessToken})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return &wrappers.BoolValue{Value: false}, nil
			}
		}
		return nil, grpcerr.E(codes.Internal, "session cannot be fetched", err)
	}
	uid, err := session.ActorID(ses.Session.SubjectId).UserID()
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "invalid actor id", err)
	}
	exists, err := iah.repository.user.Exists(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return &wrappers.BoolValue{Value: false}, nil
		}
		return nil, grpcerr.E(codes.Internal, "user cannot be fetched", err)
	}

	return &wrappers.BoolValue{Value: exists}, nil
}
