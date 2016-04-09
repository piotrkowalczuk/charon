package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listUserGroupsHandler struct {
	*handler
}

// TODO: missing firewall
func (lugh *listUserGroupsHandler) handle(ctx context.Context, req *ListUserGroupsRequest) (*ListUserGroupsResponse, error) {
	lugh.loggerWith("user_id", req.Id)

	ents, err := lugh.repository.group.FindByUserID(req.Id)
	if err != nil {
		return nil, err
	}

	return lugh.response(ents)
}

func (lugh *listUserGroupsHandler) firewall(req *ListUserGroupsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(UserGroupCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charon: list of user groups cannot be retrieved, missing permission")
}

func (lugh *listUserGroupsHandler) response(ents []*groupEntity) (*ListUserGroupsResponse, error) {
	resp := &ListUserGroupsResponse{
		Groups: make([]*Group, 0, len(ents)),
	}
	var (
		err error
		msg *Group
	)
	for _, e := range ents {
		if msg, err = e.Message(); err != nil {
			return nil, err
		}
		resp.Groups = append(resp.Groups, msg)
	}

	return resp, nil
}
