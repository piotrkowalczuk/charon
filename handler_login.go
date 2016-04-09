package charon

import (
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type loginHandler struct {
	*handler
	hasher PasswordHasher
}

func (lh *loginHandler) handle(ctx context.Context, r *LoginRequest) (*LoginResponse, error) {
	lh.logger = log.NewContext(lh.logger).With("username", r.Username)

	if r.Username == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "charon: empty username")
	}
	if len(r.Password) == 0 {
		return nil, grpc.Errorf(codes.Unauthenticated, "charon: empty password")
	}

	user, err := lh.repository.user.FindOneByUsername(r.Username)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "charon: the username and password do not match")
	}

	if matches := lh.hasher.Compare(user.Password, []byte(r.Password)); !matches {
		return nil, grpc.Errorf(codes.Unauthenticated, "charon: the username and password do not match")
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
		return nil, grpc.Errorf(codes.Unauthenticated, "charon: user is not confirmed")
	}

	if !user.IsActive {
		return nil, grpc.Errorf(codes.Unauthenticated, "charon: user is not active")
	}

	session, err := lh.session.Start(ctx, SubjectIDFromInt64(user.ID).String(), map[string]string{
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
		return nil, grpc.Errorf(codes.Internal, "charon: last login update failure: %s", err)
	}

	return &LoginResponse{AccessToken: session.AccessToken}, nil
}
