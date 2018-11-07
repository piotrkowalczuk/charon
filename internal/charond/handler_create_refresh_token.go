package charond

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/refreshtoken"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

type createRefreshTokenHandler struct {
	*handler
}

func (crth *createRefreshTokenHandler) Create(ctx context.Context, req *charonrpc.CreateRefreshTokenRequest) (*charonrpc.CreateRefreshTokenResponse, error) {
	var (
		expireAt time.Time
		err      error
	)
	if req.ExpireAt != nil {
		if expireAt, err = ptypes.Timestamp(req.ExpireAt); err != nil {
			return nil, grpcerr.E(codes.InvalidArgument, "invalid format of expire at", err)
		}
	}

	act, err := crth.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = crth.firewall(req, act); err != nil {
		return nil, err
	}

	tkn, err := refreshtoken.Random()
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "refresh token generation failure", err)
	}

	ent, err := crth.repository.refreshToken.Create(ctx, &model.RefreshTokenEntity{
		UserID: act.User.ID,
		Token:  tkn,
		ExpireAt: pq.NullTime{
			Time:  expireAt.UTC(),
			Valid: !expireAt.IsZero(),
		},
		CreatedBy: ntypes.Int64{Int64: act.User.ID, Valid: true},
		Notes:     allocNilString(req.Notes),
	})
	if err != nil {
		switch model.ErrorConstraint(err) {
		case model.TableRefreshTokenConstraintCreatedByForeignKey:
			return nil, grpcerr.E(codes.NotFound, "such user does not exist")
		case model.TableRefreshTokenConstraintTokenUnique:
			return nil, grpcerr.E(codes.AlreadyExists, "such refresh token already exists")
		case model.TableRefreshTokenConstraintUserIDForeignKey:
			return nil, grpcerr.E(codes.NotFound, "such user does not exist")
		default:
			return nil, grpcerr.E(codes.Internal, "refresh token persistence failure", err)
		}
	}

	return crth.response(ent)
}

func (crth *createRefreshTokenHandler) firewall(req *charonrpc.CreateRefreshTokenRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.RefreshTokenCanCreate) {
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "refresh token cannot be created, missing permission")
}

func (crth *createRefreshTokenHandler) response(ent *model.RefreshTokenEntity) (*charonrpc.CreateRefreshTokenResponse, error) {
	msg, err := mapping.ReverseRefreshToken(ent)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "refresh token entity mapping failure", err)
	}
	return &charonrpc.CreateRefreshTokenResponse{
		RefreshToken: msg,
	}, nil
}
