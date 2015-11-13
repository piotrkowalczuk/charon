package main

import (
	"strconv"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Login ...
func (rs *rpcServer) Login(ctx context.Context, r *charon.LoginRequest) (*charon.LoginResponse, error) {
	user, err := rs.userRepository.FindOneByUsername(r.Username)
	if err != nil {
		sklog.Debug(rs.logger, "login failed, user with such username does not exists", "username", r.Username)

		return nil, grpc.Errorf(codes.Unauthenticated, "charond: the username and password do not match")
	}

	if matches := rs.passwordHasher.Compare(user.Password, r.Password); !matches {
		sklog.Debug(rs.logger, "login failed, wrong password", "username", r.Username)

		return nil, grpc.Errorf(codes.Unauthenticated, "charond: the username and password do not match")
	}

	if !user.IsConfirmed {
		sklog.Debug(rs.logger, "login failed, email confirmation is missing", r.Username)

		return nil, grpc.Errorf(codes.Unauthenticated, "charond: user is not confirmed")
	}

	if !user.IsActive {
		sklog.Debug(rs.logger, "login failed, user is not active", r.Username)

		return nil, grpc.Errorf(codes.Unauthenticated, "charond: user is not active")
	}

	resp, err := rs.mnemosyne.Create(context.Background(), &mnemosyne.CreateRequest{
		Data: map[string]string{
			"user_id":    strconv.FormatInt(user.ID, 10),
			"username":   user.Username,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
	if err != nil {
		sklog.Error(rs.logger, err, "username", r.Username)

		return nil, err
	}

	err = rs.userRepository.UpdateLastLoginAt(user.ID)
	if err != nil {
		sklog.Error(rs.logger, err, "username", r.Username)

		return nil, grpc.Errorf(codes.Internal, "charond: last login update failure: %s", err)
	}

	return &charon.LoginResponse{Session: resp.Session}, nil
}

// Logout ...
func (rs *rpcServer) Logout(ctx context.Context, r *charon.LogoutRequest) (*charon.LogoutResponse, error) {
	if r.Token.String() == "" { // TODO: probably wrong, implement IsEmpty method for ID
		return nil, grpc.Errorf(codes.InvalidArgument, "charond: empty session id, logout aborted")
	}

	resp, err := rs.mnemosyne.Abandon(context.Background(), &mnemosyne.AbandonRequest{Token: r.Token})
	if err != nil {
		sklog.Error(rs.logger, err, "session_id", r.Token)

		return nil, err
	}

	if !resp.Abandoned {
		sklog.Debug(rs.logger, "mnemosyne responded without error but session was not abandoned, propably does not exists", "session_id", r.Token)
	} else {
		sklog.Debug(rs.logger, "successful logout", "session_id", r.Token)
	}

	return &charon.LogoutResponse{}, nil
}

// IsGranted ...
func (rs *rpcServer) IsGranted(ctx context.Context, r *charon.IsGrantedRequest) (*charon.IsGrantedResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "is granted is not implemented yet")
}

// HasPrivileges ...
func (rs *rpcServer) HasPrivileges(ctx context.Context, r *charon.HasPrivilegesRequest) (*charon.HasPrivilegesResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "has privileges is not implemented yet")
}
