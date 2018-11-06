package charond

import (
	"bytes"
	"context"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/sklog"

	"google.golang.org/grpc/codes"
)

type loginHandler struct {
	*handler
	hasher password.Hasher
}

func (lh *loginHandler) Login(ctx context.Context, r *charonrpc.LoginRequest) (*wrappers.StringValue, error) {
	if r.Username == "" {
		return nil, grpcerr.E(codes.InvalidArgument, "empty username")
	}
	if len(r.Password) == 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "empty password")
	}

	var (
		err error
		usr *model.UserEntity
	)

	usr, err = lh.repository.user.FindOneByUsername(ctx, r.Username)
	if err != nil {
		return nil, grpcerr.E(codes.Unauthenticated, "user does not exists")
	}

	if bytes.Equal(usr.Password, model.ExternalPassword) {
		return nil, grpcerr.E(codes.FailedPrecondition, "authentication failure, external password manager not implemented")
	}
	if matches := lh.hasher.Compare(usr.Password, []byte(r.Password)); !matches {
		return nil, grpcerr.E(codes.Unauthenticated, "the username and password do not match")
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

	sklog.Debug(lh.logger, "user session has been started", "user_id", usr.ID)

	_, err = lh.repository.user.UpdateLastLoginAt(ctx, usr.ID)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "last login update failure: %s", err)
	}

	sklog.Debug(lh.logger, "user last login at field has been updated", "user_id", usr.ID)

	return &wrappers.StringValue{Value: res.Session.AccessToken}, nil
}
