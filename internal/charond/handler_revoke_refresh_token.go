package charond

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/ntypes"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type revokeRefreshTokenHandler struct {
	*handler
}

func (h *revokeRefreshTokenHandler) Revoke(ctx context.Context, req *charonrpc.RevokeRefreshTokenRequest) (*charonrpc.RevokeRefreshTokenResponse, error) {
	if len(req.Token) == 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "refresh token cannot be disabled, invalid token: %s", req.Token)
	}
	if req.UserId == 0 {
		return nil, grpcerr.E(codes.InvalidArgument, "refresh token cannot be disabled, missing user id")
	}

	act, err := h.Actor(ctx)
	if err != nil {
		return nil, err
	}
	ent, err := h.repository.refreshToken.FindOneByTokenAndUserID(ctx, req.Token, req.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpcerr.E(codes.NotFound, "refresh token does not exists")
		}
		return nil, grpcerr.E(codes.Internal, "refresh token could not be retrieved", err)
	}
	if err = h.firewall(req, act, ent); err != nil {
		return nil, err
	}

	ent, err = h.repository.refreshToken.UpdateOneByToken(ctx, req.Token, &model.RefreshTokenPatch{
		Revoked: ntypes.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpcerr.E(codes.NotFound, "refresh token does not exists")
		}
		return nil, grpcerr.E(codes.Internal, "refresh token could not be disabled", err)
	}

	res, err := h.session.Delete(ctx, &mnemosynerpc.DeleteRequest{
		SubjectId:    session.ActorIDFromInt64(ent.UserID).String(),
		RefreshToken: req.Token,
	})
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return nil, grpcerr.E(codes.Internal, "session could not be removed", err)
		}
	}
	h.logger.Debug("refresh token corresponding sessions removed", zap.Int64("count", res.Value))

	msg, err := mapping.ReverseRefreshToken(ent)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "refresh token mapping failure", err)
	}
	return &charonrpc.RevokeRefreshTokenResponse{
		RefreshToken: msg,
	}, nil
}

func (h *revokeRefreshTokenHandler) firewall(req *charonrpc.RevokeRefreshTokenRequest, act *session.Actor, ent *model.RefreshTokenEntity) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.RefreshTokenCanRevokeAsStranger) {
		return nil
	}
	if act.Permissions.Contains(charon.RefreshTokenCanRevokeAsOwner) {
		if act.User.ID == ent.UserID {
			return nil
		}
		return grpcerr.E(codes.PermissionDenied, "refresh token cannot be revoked by stranger, missing permission")
	}
	return grpcerr.E(codes.PermissionDenied, "refresh token cannot be revoked, missing permission")
}
