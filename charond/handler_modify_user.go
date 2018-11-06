package charond

import (
	"context"
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type modifyUserHandler struct {
	*handler
}

func (muh *modifyUserHandler) Modify(ctx context.Context, req *charonrpc.ModifyUserRequest) (*charonrpc.ModifyUserResponse, error) {
	if req.Id <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "user cannot be modified, invalid ID: %d", req.Id)
	}

	act, err := muh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	ent, err := muh.repository.user.FindOneByID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "user does not exists")
		}
		return nil, err
	}

	if hint, ok := muh.firewall(req, ent, act); !ok {
		return nil, grpc.Errorf(codes.PermissionDenied, hint)
	}

	ent, err = muh.repository.user.UpdateOneByID(ctx, req.Id, &model.UserPatch{
		FirstName:   allocNilString(req.FirstName),
		IsActive:    allocNilBool(req.IsActive),
		IsConfirmed: allocNilBool(req.IsConfirmed),
		IsStaff:     allocNilBool(req.IsStaff),
		IsSuperuser: allocNilBool(req.IsSuperuser),
		LastName:    allocNilString(req.LastName),
		Password:    req.SecurePassword,
		UpdatedBy:   ntypes.Int64{Int64: act.User.ID, Valid: act.User.ID != 0},
		Username:    allocNilString(req.Username),
	})
	if err != nil {
		switch model.ErrorConstraint(err) {
		case model.TableUserConstraintUsernameUnique:
			return nil, grpc.Errorf(codes.AlreadyExists, "user with such username already exists")
		default:
			return nil, err
		}
	}

	return muh.response(ent)
}

func (muh *modifyUserHandler) firewall(req *charonrpc.ModifyUserRequest, ent *model.UserEntity, act *session.Actor) (string, bool) {
	isOwner := act.User.ID == ent.ID

	if !act.User.IsSuperuser {
		switch {
		case ent.IsSuperuser:
			return "only superuser can modify a superuser account", false
		case ent.IsStaff && !isOwner && !act.Permissions.Contains(charon.UserCanModifyStaffAsStranger):
			return "missing permission to modify staff account as a stranger", false
		case ent.IsStaff && isOwner && !act.Permissions.Contains(charon.UserCanModifyStaffAsOwner):
			return "missing permission to modify staff account as an owner", false
		case req.IsSuperuser != nil && req.IsSuperuser.Valid:
			return "only superuser can change existing account to superuser", false
		case req.IsStaff != nil && req.IsStaff.Valid && !act.Permissions.Contains(charon.UserCanCreateStaff):
			return "missing permission to change existing account to staff", false
		}
	}

	return "", true
}

func (muh *modifyUserHandler) response(u *model.UserEntity) (*charonrpc.ModifyUserResponse, error) {
	msg, err := u.Message()
	if err != nil {
		return nil, err
	}
	return &charonrpc.ModifyUserResponse{
		User: msg,
	}, nil
}
