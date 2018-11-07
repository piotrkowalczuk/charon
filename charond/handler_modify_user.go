package charond

import (
	"context"
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"

	"google.golang.org/grpc/codes"
)

type modifyUserHandler struct {
	*handler
}

func (muh *modifyUserHandler) Modify(ctx context.Context, req *charonrpc.ModifyUserRequest) (*charonrpc.ModifyUserResponse, error) {
	if req.Id <= 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "user cannot be modified, invalid id")
	}

	act, err := muh.Actor(ctx)
	if err != nil {
		return nil, err
	}

	if !muh.firewall(req, act) {
		return nil, grpcerr.E(codes.PermissionDenied, "user cannot be modified, missing permissions")
	}

	ent, err := muh.repository.user.FindOneByID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpcerr.E(codes.NotFound, "user does not exists")
		}
		return nil, grpcerr.E(codes.Internal, "find user by id query failed", err)
	}

	if err := muh.firewallEntity(req, ent, act); err != nil {
		return nil, err
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
			return nil, grpcerr.E(codes.AlreadyExists, "user with such username already exists")
		default:
			if err == sql.ErrNoRows {
				return nil, grpcerr.E(codes.NotFound, "user does not exists")
			}
			return nil, grpcerr.E(codes.Internal, "update user by id query failed", err)
		}
	}

	return muh.response(ent)
}

func (muh *modifyUserHandler) firewall(req *charonrpc.ModifyUserRequest, act *session.Actor) bool {
	return act.User.IsSuperuser || act.Permissions.Contains(
		charon.UserCanModifyAsStranger,
		charon.UserCanModifyAsOwner,
		charon.UserCanModifyStaffAsStranger,
		charon.UserCanModifyStaffAsOwner,
	)
}

func (muh *modifyUserHandler) firewallEntity(req *charonrpc.ModifyUserRequest, ent *model.UserEntity, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if ent.IsSuperuser {
		return grpcerr.E(codes.PermissionDenied, "only superuser can modify another superuser")
	}
	if req.IsSuperuser.BoolOr(false) {
		return grpcerr.E(codes.PermissionDenied, "only superuser can promote another user to become superuser")
	}
	// STAFF USERS
	if ent.IsStaff {
		if ent.CreatedBy.Int64Or(0) == act.User.ID {
			if !act.Permissions.Contains(charon.UserCanModifyStaffAsStranger, charon.UserCanModifyStaffAsOwner) {
				return grpcerr.E(codes.PermissionDenied, "staff user cannot be modified as an owner, missing permission")
			}
			return nil
		}
		if !act.Permissions.Contains(charon.UserCanModifyStaffAsStranger) {
			return grpcerr.E(codes.PermissionDenied, "staff user cannot be modified as an stranger, missing permission")
		}
		return nil
	}
	if req.IsStaff.BoolOr(false) {
		if !act.Permissions.Contains(charon.UserCanCreateStaff) {
			return grpcerr.E(codes.PermissionDenied, "regular user cannot be promoted to staff, missing permission")
		}
	}
	// NON STAFF USERS
	if ent.CreatedBy.Int64Or(0) == act.User.ID {
		if !act.Permissions.Contains(charon.UserCanModifyAsStranger, charon.UserCanModifyAsOwner) {
			return grpcerr.E(codes.PermissionDenied, "user cannot be modified as an owner, missing permission")
		}
		return nil
	}
	if !act.Permissions.Contains(charon.UserCanModifyAsStranger) {
		return grpcerr.E(codes.PermissionDenied, "user cannot be modified as a stranger, missing permission")
	}
	return nil
}

func (muh *modifyUserHandler) response(ent *model.UserEntity) (*charonrpc.ModifyUserResponse, error) {
	msg, err := mapping.ReverseUser(ent)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "user reverse mapping failure")
	}
	return &charonrpc.ModifyUserResponse{
		User: msg,
	}, nil
}
