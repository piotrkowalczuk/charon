package charond

import (
	"github.com/piotrkowalczuk/charon"
	charonrpc "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"
	"github.com/piotrkowalczuk/charon/internal/grpcerr"
	"github.com/piotrkowalczuk/charon/internal/mapping"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/qtypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

type listRefreshTokensHandler struct {
	*handler
}

func (lrth *listRefreshTokensHandler) List(ctx context.Context, req *charonrpc.ListRefreshTokensRequest) (*charonrpc.ListRefreshTokensResponse, error) {
	act, err := lrth.Actor(ctx)
	if err != nil {
		return nil, err
	}
	if err = lrth.firewall(req, act); err != nil {
		return nil, err
	}

	ents, err := lrth.repository.refreshToken.Find(ctx, &model.RefreshTokenFindExpr{
		Limit:   req.GetLimit().Int64Or(10),
		Offset:  req.GetOffset().Int64Or(0),
		OrderBy: mapping.OrderBy(req.GetOrderBy()),
		Where:   mapping.RefreshTokenQuery(req.GetQuery()),
	})
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "find refresh token query failed", err)
	}

	msg, err := mapping.ReverseRefreshTokens(ents)
	if err != nil {
		return nil, grpcerr.E(codes.Internal, "refresh token reverse mapping failure")
	}
	return &charonrpc.ListRefreshTokensResponse{
		RefreshTokens: msg,
	}, nil
}

func (lrth *listRefreshTokensHandler) firewall(req *charonrpc.ListRefreshTokensRequest, act *session.Actor) error {
	if act.User.IsSuperuser {
		return nil
	}
	if act.Permissions.Contains(charon.RefreshTokenCanRetrieveAsStranger) {
		return nil
	}
	if act.Permissions.Contains(charon.RefreshTokenCanRetrieveAsOwner) {
		if req.Query == nil {
			req.Query = &charonrpc.RefreshTokenQuery{}
		}
		req.Query.UserId = qtypes.EqualInt64(act.User.ID)
		return nil
	}

	return grpcerr.E(codes.PermissionDenied, "list of refresh tokens cannot be retrieved, missing permission")
}
