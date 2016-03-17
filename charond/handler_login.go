package main

import (
	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type loginHandler struct {
	*handler
	hasher charon.PasswordHasher
}

func (lh *loginHandler) handle(ctx context.Context, r *charon.LoginRequest) (*charon.LoginResponse, error) {
	lh.logger = log.NewContext(lh.logger).With("username", r.Username)

	if r.Username == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "charond: empty username")
	}
	if len(r.Password) == 0 {
		return nil, grpc.Errorf(codes.Unauthenticated, "charond: empty password")
	}

	user, err := lh.repository.user.FindOneByUsername(r.Username)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "charond: the username and password do not match")
	}

	if matches := lh.hasher.Compare(user.Password, []byte(r.Password)); !matches {
		return nil, grpc.Errorf(codes.Unauthenticated, "charond: the username and password do not match")
	}

	lh.loggerWith(
		"is_confirmed", user.IsConfirmed,
		"is_staff", user.IsStaff,
		"is_superuser", user.IsSuperuser,
		"is_active", user.IsActive,
		"first_name", user.FirstName,
		"last_name", user.LastName,
	)
	if !user.IsConfirmed {
		return nil, grpc.Errorf(codes.Unauthenticated, "charond: user is not confirmed")
	}

	if !user.IsActive {
		return nil, grpc.Errorf(codes.Unauthenticated, "charond: user is not active")
	}

	session, err := lh.session.Start(ctx, charon.SubjectIDFromInt64(user.ID).String(), map[string]string{
		"username":   user.Username,
		"filht_name": user.FirstName,
		"last_name":  user.LastName,
	})
	if err != nil {
		return nil, err
	}

	lh.loggerWith("token", session.AccessToken.Encode())

	_, err = lh.repository.user.UpdateLastLoginAt(user.ID)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "charond: last login update failure: %s", err)
	}

	return &charon.LoginResponse{AccessToken: session.AccessToken}, nil
}
