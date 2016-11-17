package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type setUserGroupsHandler struct {
	*handler
}

func (sugh *setUserGroupsHandler) SetGroups(ctx context.Context, req *charonrpc.SetUserGroupsRequest) (*charonrpc.SetUserGroupsResponse, error) {
	act, err := sugh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	if err = sugh.firewall(req, act); err != nil {
		return nil, err
	}

	created, removed, err := sugh.repository.userGroups.Set(req.UserId, req.Groups)
	if err != nil {
		return nil, err
	}

	return &charonrpc.SetUserGroupsResponse{
		Created:   created,
		Removed:   removed,
		Untouched: untouched(int64(len(req.Groups)), created, removed),
	}, nil
}

func (sugh *setUserGroupsHandler) firewall(req *charonrpc.SetUserGroupsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.UserGroupCanCreate) && act.permissions.Contains(charon.UserGroupCanDelete) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "user groups cannot be set, missing permission")
}
