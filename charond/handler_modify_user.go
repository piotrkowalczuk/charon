package main

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/nilt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type modifyUserHandler struct {
	*handler
}

func (muh *modifyUserHandler) handle(ctx context.Context, req *charon.ModifyUserRequest) (*charon.ModifyUserResponse, error) {
	if req.Id <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "charond: user cannot be modified, invalid id: %d", req.Id)
	}

	muh.loggerWith("user_id", req.Id)

	actor, err := muh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	entity, err := muh.repository.user.FindOneByID(req.Id)
	if err != nil {
		return nil, err
	}

	if hint, ok := muh.firewall(req, entity, actor); !ok {
		return nil, grpc.Errorf(codes.PermissionDenied, "charond: "+hint)
	}

	entity, err = muh.repository.user.UpdateByID(
		req.Id,
		nil,
		nil,
		nilt.Int64{Int64: actor.user.ID, Valid: actor.user.ID != 0},
		nilString(req.FirstName),
		nilBool(req.IsActive),
		nilBool(req.IsConfirmed),
		nilBool(req.IsStaff),
		nilBool(req.IsSuperuser),
		nil,
		nilString(req.LastName),
		req.SecurePassword,
		nil,
		nilt.Int64{},
		nilString(req.Username),
	)
	if err != nil {
		return nil, mapUserError(err)
	}

	return &charon.ModifyUserResponse{
		User: entity.Message(),
	}, nil
}

func (muh *modifyUserHandler) firewall(req *charon.ModifyUserRequest, entity *userEntity, actor *actor) (string, bool) {
	isOwner := actor.user.ID == entity.ID

	if !actor.user.IsSuperuser {
		switch {
		case entity.IsSuperuser:
			return "only superuser can modify a superuser account", false
		case entity.IsStaff && !isOwner && actor.permissions.Contains(charon.UserCanModifyStaffAsStranger):
			return "missing permission to modify an account as a stranger", false
		case entity.IsStaff && isOwner && actor.permissions.Contains(charon.UserCanModifyStaffAsOwner):
			return "missing permission to modify an account as an owner", false
		case req.IsSuperuser != nil && req.IsSuperuser.Valid:
			return "only superuser can change existing account to superuser", false
		case req.IsStaff != nil && req.IsStaff.Valid && !actor.permissions.Contains(charon.UserCanCreateStaff):
			return "user is not allowed to create user with is_staff property that has custom value", false
		}
	}

	return "", true
}
