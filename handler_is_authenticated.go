package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type isAuthenticatedHandler struct {
	*handler
}

func (iah *isAuthenticatedHandler) handle(ctx context.Context, req *IsAuthenticatedRequest) (*IsAuthenticatedResponse, error) {
	if req.AccessToken == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "charon: authentication status cannot be checked, missing token")
	}

	iah.loggerWith("token", req.AccessToken.Encode())

	ses, err := iah.session.Get(ctx, *req.AccessToken)
	if err != nil {
		return nil, err
	}
	uid, err := SubjectID(ses.SubjectId).UserID()
	if err != nil {
		return nil, err
	}
	iah.loggerWith("user_id", uid)
	exists, err := iah.repository.user.Exists(uid)
	if err != nil {
		return nil, err
	}

	return &IsAuthenticatedResponse{
		Authenticated: exists,
	}, nil
}
