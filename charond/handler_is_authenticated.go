package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type isAuthenticatedHandler struct {
	*handler
}

func (iah *isAuthenticatedHandler) handle(ctx context.Context, req *charon.IsAuthenticatedRequest) (*charon.IsAuthenticatedResponse, error) {
	if req.AccessToken == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "authentication status cannot be checked, missing token")
	}

	iah.loggerWith("token", req.AccessToken)

	ses, err := iah.session.Get(ctx, &mnemosynerpc.GetRequest{AccessToken: req.AccessToken})
	if err != nil {
		return nil, handleMnemosyneError(err)
	}
	uid, err := charon.SubjectID(ses.Session.SubjectId).UserID()
	if err != nil {
		return nil, err
	}
	iah.loggerWith("user_id", uid)
	exists, err := iah.repository.user.Exists(uid)
	if err != nil {
		return nil, err
	}

	return &charon.IsAuthenticatedResponse{
		Authenticated: exists,
	}, nil
}
