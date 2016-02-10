package main

import (
	"github.com/pborman/uuid"
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type createUserHandler struct {
	*handler
	hasher charon.PasswordHasher
}

func (cuh *createUserHandler) handle(ctx context.Context, req *charon.CreateUserRequest) (*charon.CreateUserResponse, error) {
	cuh.loggerWith("username", req.Username, "superuser", req.IsSuperuser.BoolOr(false))

	akt, err := cuh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = cuh.firewall(req, akt); err != nil {
		return nil, err
	}

	if len(req.SecurePassword) == 0 {
		req.SecurePassword, err = cuh.hasher.Hash([]byte(req.PlainPassword))
		if err != nil {
			return nil, err
		}
	} else {
		if !akt.user.IsSuperuser {
			return nil, grpc.Errorf(codes.PermissionDenied, "charond: only superuser can create an user with manualy defined secure password")
		}
	}

	entity, err := cuh.repository.user.Create(
		req.Username,
		req.SecurePassword,
		req.FirstName,
		req.LastName,
		uuid.NewRandom(),
		req.IsSuperuser.BoolOr(false),
		req.IsStaff.BoolOr(false),
		req.IsActive.BoolOr(false),
		req.IsConfirmed.BoolOr(false),
	)
	if err != nil {
		return nil, mapUserError(err)
	}

	return &charon.CreateUserResponse{
		User: entity.Message(),
	}, nil
}

func (cuh *createUserHandler) firewall(req *charon.CreateUserRequest, akt *actor) error {
	if akt.isLocalhost() || akt.user.IsSuperuser {
		return nil
	}
	if req.IsSuperuser.BoolOr(false) {
		return grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create superuser")
	}
	if req.IsStaff.BoolOr(false) && !akt.permissions.Contains(charon.UserCanCreateStaff) {
		return grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create staff user")
	}
	if !akt.permissions.Contains(charon.UserCanCreateStaff, charon.UserCanCreate) {
		return grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create another user")
	}

	return nil
}
