package charon

import (
	"github.com/piotrkowalczuk/nilt"
	"github.com/piotrkowalczuk/pqt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type modifyUserHandler struct {
	*handler
}

func (muh *modifyUserHandler) handle(ctx context.Context, req *ModifyUserRequest) (*ModifyUserResponse, error) {
	if req.Id <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "charon: user cannot be modified, invalid id: %d", req.Id)
	}

	muh.loggerWith("user_id", req.Id)

	actor, err := muh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	ent, err := muh.repository.user.FindOneByID(req.Id)
	if err != nil {
		return nil, err
	}

	if hint, ok := muh.firewall(req, ent, actor); !ok {
		return nil, grpc.Errorf(codes.PermissionDenied, "charon: "+hint)
	}

	ent, err = muh.repository.user.UpdateByID(
		req.Id,
		nil,
		nil,
		&nilt.Int64{Int64: actor.user.ID, Valid: actor.user.ID != 0},
		req.FirstName,
		req.IsActive,
		req.IsConfirmed,
		req.IsStaff,
		req.IsSuperuser,
		nil,
		req.LastName,
		req.SecurePassword,
		nil,
		nil,
		req.Username,
	)
	if err != nil {
		switch pqt.ErrorConstraint(err) {
		case tableUserConstraintUsernameUnique:
			return nil, grpc.Errorf(codes.AlreadyExists, ErrDescUserWithUsernameExists)
		default:
			return nil, err
		}
	}

	return muh.response(ent)
}

func (muh *modifyUserHandler) firewall(req *ModifyUserRequest, entity *userEntity, actor *actor) (string, bool) {
	isOwner := actor.user.ID == entity.ID

	if !actor.user.IsSuperuser {
		switch {
		case entity.IsSuperuser:
			return "only superuser can modify a superuser account", false
		case entity.IsStaff && !isOwner && actor.permissions.Contains(UserCanModifyStaffAsStranger):
			return "missing permission to modify an account as a stranger", false
		case entity.IsStaff && isOwner && actor.permissions.Contains(UserCanModifyStaffAsOwner):
			return "missing permission to modify an account as an owner", false
		case req.IsSuperuser != nil && req.IsSuperuser.Valid:
			return "only superuser can change existing account to superuser", false
		case req.IsStaff != nil && req.IsStaff.Valid && !actor.permissions.Contains(UserCanCreateStaff):
			return "user is not allowed to create user with is_staff property that has custom value", false
		}
	}

	return "", true
}

func (muh *modifyUserHandler) response(u *userEntity) (*ModifyUserResponse, error) {
	msg, err := u.Message()
	if err != nil {
		return nil, err
	}
	return &ModifyUserResponse{
		User: msg,
	}, nil
}
