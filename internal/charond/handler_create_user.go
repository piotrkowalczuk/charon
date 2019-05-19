package charond

import (
	"context"

	"github.com/google/uuid"
	"github.com/piotrkowalczuk/charon"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/password"
	"github.com/piotrkowalczuk/charon/internal/session"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type createUserHandler struct {
	*handler
	hasher password.Hasher
}

func (cuh *createUserHandler) Create(ctx context.Context, req *charonrpc.CreateUserRequest) (*charonrpc.CreateUserResponse, error) {
	if len(req.Username) < 3 {
		return nil, grpcerr.E(codes.InvalidArgument, "username needs to be at least 3 characters long")
	}
	if len(req.SecurePassword) == 0 {
		if len(req.PlainPassword) < 8 {
			return nil, grpcerr.E(codes.InvalidArgument, "password needs to be at least 8 characters long")
		}
	}

	act, err := cuh.Actor(ctx)
	if err != nil {
		if req.IsSuperuser.BoolOr(false) {
			count, err := cuh.repository.user.Count(ctx)
			if err != nil {
				return nil, grpcerr.E(codes.Internal, "number of users cannot be checked", err)
			}
			if count > 0 {
				return nil, grpcerr.E(codes.AlreadyExists, "initial superuser account already exists")
			}

			// If session.Actor does not exists, even single user does not exists and request contains IsSuperuser equals to true.
			// Then move forward, its request that is trying to create first user (that needs to be a superuser).
		} else {
			return nil, err
		}
	} else {
		if err = cuh.firewall(req, act); err != nil {
			return nil, err
		}
	}

	if len(req.SecurePassword) == 0 {
		req.SecurePassword, err = cuh.hasher.Hash([]byte(req.PlainPassword))
		if err != nil {
			return nil, grpcerr.E(codes.Internal, "password hashing failure", err)
		}
	} else {
		if !act.User.IsSuperuser {
			return nil, grpcerr.E(codes.PermissionDenied, "only superuser can create an user with manually defined secure password")
		}
	}

	token, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "confirmation token generation failure: %s", err)
	}
	ent, err := cuh.repository.user.Create(ctx, &model.UserEntity{
		Username:          req.Username,
		Password:          req.SecurePassword,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		ConfirmationToken: token[:],
		IsSuperuser:       req.IsSuperuser.BoolOr(false),
		IsStaff:           req.IsStaff.BoolOr(false),
		IsActive:          req.IsActive.BoolOr(false),
		IsConfirmed:       req.IsConfirmed.BoolOr(false),
	})
	if err != nil {
		switch model.ErrorConstraint(err) {
		case model.TableUserConstraintUsernameUnique:
			return nil, grpcerr.E(codes.AlreadyExists, "user with such username already exists")
		default:
			return nil, grpcerr.E(codes.Internal, "user cannot be persisted", err)
		}
	}

	return cuh.response(ent)
}

func (cuh *createUserHandler) firewall(req *charonrpc.CreateUserRequest, act *session.Actor) error {
	if act.IsLocal || act.User.IsSuperuser {
		return nil
	}
	if req.IsSuperuser.BoolOr(false) {
		return grpcerr.E(codes.PermissionDenied, "user is not allowed to create superuser")
	}
	if req.IsStaff.BoolOr(false) && !act.Permissions.Contains(charon.UserCanCreateStaff) {
		return grpcerr.E(codes.PermissionDenied, "user is not allowed to create staff user")
	}
	if !act.Permissions.Contains(charon.UserCanCreateStaff, charon.UserCanCreate) {
		return grpcerr.E(codes.PermissionDenied, "user is not allowed to create another user")
	}

	return nil
}

func (cuh *createUserHandler) response(ent *model.UserEntity) (*charonrpc.CreateUserResponse, error) {
	msg, err := mapping.ReverseUser(ent)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "user entity mapping failure", err)
	}
	return &charonrpc.CreateUserResponse{User: msg}, nil
}
