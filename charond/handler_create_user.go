package charond

import (
	"github.com/pborman/uuid"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/pqt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type createUserHandler struct {
	*handler
	hasher charon.PasswordHasher
}

func (cuh *createUserHandler) handle(ctx context.Context, req *charon.CreateUserRequest) (*charon.CreateUserResponse, error) {
	cuh.loggerWith("username", req.Username, "is_superuser", req.IsSuperuser.BoolOr(false))

	act, err := cuh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = cuh.firewall(req, act); err != nil {
		return nil, err
	}

	if act.isLocal && req.IsSuperuser.BoolOr(false) {
		count, err := cuh.repository.user.Count()
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, grpc.Errorf(codes.AlreadyExists, "charond: initial superuser account already exists")
		}
	}
	if len(req.SecurePassword) == 0 {
		req.SecurePassword, err = cuh.hasher.Hash([]byte(req.PlainPassword))
		if err != nil {
			return nil, err
		}
	} else {
		// TODO: only one superuser can be defined so this else statement makes no sense in this place.
		if !act.user.IsSuperuser {
			return nil, grpc.Errorf(codes.PermissionDenied, "charond: only superuser can create an user with manualy defined secure password")
		}
	}

	ent, err := cuh.repository.user.Create(
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
		switch pqt.ErrorConstraint(err) {
		case tableUserConstraintPrimaryKey:
			return nil, grpc.Errorf(codes.AlreadyExists, charon.ErrDescUserWithIDExists)
		case tableUserConstraintUsernameUnique:
			return nil, grpc.Errorf(codes.AlreadyExists, charon.ErrDescUserWithUsernameExists)
		default:
			return nil, err
		}
	}

	return cuh.response(ent)
}

func (cuh *createUserHandler) firewall(req *charon.CreateUserRequest, act *actor) error {
	if act.isLocal || act.user.IsSuperuser {
		return nil
	}
	if req.IsSuperuser.BoolOr(false) {
		return grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create superuser")
	}
	if req.IsStaff.BoolOr(false) && !act.permissions.Contains(charon.UserCanCreateStaff) {
		return grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create staff user")
	}
	if !act.permissions.Contains(charon.UserCanCreateStaff, charon.UserCanCreate) {
		return grpc.Errorf(codes.PermissionDenied, "charond: user is not allowed to create another user")
	}

	return nil
}

func (cuh *createUserHandler) response(ent *userEntity) (*charon.CreateUserResponse, error) {
	msg, err := ent.message()
	if err != nil {
		return nil, err
	}
	return &charon.CreateUserResponse{
		User: msg,
	}, nil
}
