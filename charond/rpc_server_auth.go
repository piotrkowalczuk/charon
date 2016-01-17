package main

import (
	"fmt"

	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// Login implements charon.RPCServer interface.
func (rs *rpcServer) Login(ctx context.Context, r *charon.LoginRequest) (*charon.LoginResponse, error) {
	if r.Username == "" {
		sklog.Debug(rs.logger, "login failed, empty username")

		return nil, grpc.Errorf(codes.Unauthenticated, "charond: empty username")
	}
	if len(r.Password) == 0 {
		sklog.Debug(rs.logger, "login failed, empty password", "username", r.Username)

		return nil, grpc.Errorf(codes.Unauthenticated, "charond: empty password")
	}

	user, err := rs.repository.user.FindOneByUsername(r.Username)
	if err != nil {
		sklog.Debug(rs.logger, "login failed, user with such username does not exists", "username", r.Username)

		return nil, grpc.Errorf(codes.Unauthenticated, "charond: the username and password do not match")
	}

	if matches := rs.passwordHasher.Compare(user.Password, []byte(r.Password)); !matches {
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

	session, err := rs.session.Start(ctx, charon.NewSessionSubjectID(user.ID).String(), map[string]string{
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})
	if err != nil {
		sklog.Error(rs.logger, err, "username", r.Username)

		return nil, err
	}

	_, err = rs.repository.user.UpdateLastLoginAt(user.ID)
	if err != nil {
		sklog.Error(rs.logger, err, "username", r.Username)

		return nil, grpc.Errorf(codes.Internal, "charond: last login update failure: %s", err)
	}

	return &charon.LoginResponse{Token: session.Token}, nil
}

// Logout implements charon.RPCServer interface.
func (rs *rpcServer) Logout(ctx context.Context, r *charon.LogoutRequest) (*charon.LogoutResponse, error) {
	if r.Token.IsEmpty() { // TODO: probably wrong, implement IsEmpty method for ID
		return nil, grpc.Errorf(codes.InvalidArgument, "charond: empty session id, logout aborted")
	}

	err := rs.session.Abandon(ctx, *r.Token)
	if err != nil {
		sklog.Error(rs.logger, err, "session_id", r.Token)

		return nil, err
	}

	sklog.Debug(rs.logger, "successful logout", "session_id", r.Token)

	return &charon.LogoutResponse{}, nil
}

// IsGranted implements charon.RPCServer interface.
func (rs *rpcServer) IsGranted(ctx context.Context, r *charon.IsGrantedRequest) (*charon.IsGrantedResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "is granted is not implemented yet")
}

// BelongsTo implements charon.RPCServer interface.
func (rs *rpcServer) BelongsTo(ctx context.Context, r *charon.BelongsToRequest) (*charon.BelongsToResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "belongs to is not implemented yet")
}

// IsAuthenticated implements charon.RPCServer interface.
func (rs *rpcServer) IsAuthenticated(ctx context.Context, r *charon.IsAuthenticatedRequest) (*charon.IsAuthenticatedResponse, error) {
	if r.Token == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "charond: authentication status cannot be checked, missing token")
	}
	ok, err := rs.session.Exists(ctx, *r.Token)
	if err != nil {
		return nil, err
	}

	return &charon.IsAuthenticatedResponse{
		Authenticated: ok,
	}, nil
}

// Subject implements charon.RPCServer interface.
func (rs *rpcServer) Subject(ctx context.Context, req *charon.SubjectRequest) (*charon.SubjectResponse, error) {
	var (
		ok bool
		md metadata.MD
	)
	if md, ok = metadata.FromContext(ctx); ok {
		grpc.Header(&md)
	}

	fmt.Println(md)

	ses, err := rs.session.FromContext(ctx)
	id, err := charon.SessionSubjectID(ses.SubjectId).UserID()
	if err != nil {
		return nil, fmt.Errorf("charond: invalid session subject id: %s", ses.SubjectId)
	}

	user, err := rs.repository.user.FindOneByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "charond: user does not exists with id: %d", id)
		}

		return nil, err
	}

	permissionEntities, err := rs.repository.permission.FindByUserID(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	permissions := make([]string, 0, len(permissionEntities))
	for _, e := range permissionEntities {
		permissions = append(permissions, e.Permission().String())
	}

	sklog.Debug(rs.logger, "subject retrieved", "subject_id", ses.SubjectId)

	return &charon.SubjectResponse{
		Id:          int64(user.ID),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Permissions: permissions,
		IsActive:    user.IsActive,
		IsConfirmed: user.IsConfirmed,
		IsStuff:     user.IsStaff,
		IsSuperuser: user.IsSuperuser,
	}, nil
}
