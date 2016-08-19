package charond

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type setUserGroupsHandler struct {
	*handler
}

func (sugh *setUserGroupsHandler) handle(ctx context.Context, req *charon.SetUserGroupsRequest) (*charon.SetUserGroupsResponse, error) {
	sugh.loggerWith("user_id", req.UserId)

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

	return &charon.SetUserGroupsResponse{
		Created:   created,
		Removed:   removed,
		Untouched: untouched(int64(len(req.Groups)), created, removed),
	}, nil
}

func (sugh *setUserGroupsHandler) firewall(req *charon.SetUserGroupsRequest, act *actor) error {
	if act.user.isSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.UserGroupCanCreate) && act.permissions.Contains(charon.UserGroupCanDelete) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "user groups cannot be set, missing permission")
}
