package charond

import (
	"context"
	"fmt"

	"github.com/piotrkowalczuk/charon/internal/service"
	"go.uber.org/zap"

	"github.com/golang/protobuf/ptypes/wrappers"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"

	"google.golang.org/grpc/codes"
)

type loginHandler struct {
	*handler
	userFinderFactory *service.UserFinderFactory
}

// Login implements charonrpc AuthClient interface.
func (lh *loginHandler) Login(ctx context.Context, r *charonrpc.LoginRequest) (*wrappers.StringValue, error) {
	if r.GetUsername() != "" || r.GetPassword() != "" {
		r.Strategy = &charonrpc.LoginRequest_UsernameAndPassword{
			UsernameAndPassword: &charonrpc.UsernameAndPasswordStrategy{
				Username: r.GetUsername(),
				Password: r.GetPassword(),
			},
		}
	}

	var (
		userFinder   service.UserFinder
		refreshToken string
	)
	switch str := r.GetStrategy().(type) {
	case *charonrpc.LoginRequest_UsernameAndPassword:
		userFinder = lh.userFinderFactory.ByUsernameAndPassword(
			str.UsernameAndPassword.GetUsername(),
			str.UsernameAndPassword.GetPassword(),
		)
	case *charonrpc.LoginRequest_RefreshToken:
		refreshToken = str.RefreshToken.GetRefreshToken()
		userFinder = lh.userFinderFactory.ByRefreshToken(refreshToken)
	}

	usr, err := userFinder.FindUser(ctx)
	if err != nil {
		fmt.Println("error", err)
		return nil, err
	}

	if !usr.IsConfirmed {
		return nil, grpcerr.E(codes.Unauthenticated, "user is not confirmed")
	}

	if !usr.IsActive {
		return nil, grpcerr.E(codes.Unauthenticated, "user is not active")
	}

	res, err := lh.session.Start(ctx, &mnemosynerpc.StartRequest{
		Session: &mnemosynerpc.Session{
			SubjectId:     session.ActorIDFromInt64(usr.ID).String(),
			SubjectClient: r.Client,
			RefreshToken:  refreshToken,
			Bag: map[string]string{
				"username":   usr.Username,
				"first_name": usr.FirstName,
				"last_name":  usr.LastName,
			},
		},
	})
	if err != nil {
		return nil, grpcerr.E("session start on login failure", err)
	}

	lh.logger.Debug("user session has been started", zap.Int64("user_id", usr.ID))

	_, err = lh.repository.user.UpdateLastLoginAt(ctx, usr.ID)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "last login update failure: %s", err)
	}

	lh.logger.Debug("user last login at field has been updated", zap.Int64("user_id", usr.ID))

	return &wrappers.StringValue{Value: res.Session.AccessToken}, nil
}
