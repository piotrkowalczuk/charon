package charond

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listUserGroupsHandler struct {
	*handler
}

// TODO: missing firewall
func (lugh *listUserGroupsHandler) handle(ctx context.Context, req *charon.ListUserGroupsRequest) (*charon.ListUserGroupsResponse, error) {
	lugh.loggerWith("user_id", req.Id)

	ents, err := lugh.repository.group.FindByUserID(req.Id)
	if err != nil {
		return nil, err
	}

	return lugh.response(ents)
}

func (lugh *listUserGroupsHandler) firewall(req *charon.ListUserGroupsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.UserGroupCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charond: list of user groups cannot be retrieved, missing permission")
}

func (lugh *listUserGroupsHandler) response(ents []*groupEntity) (*charon.ListUserGroupsResponse, error) {
	resp := &charon.ListUserGroupsResponse{
		Groups: make([]*charon.Group, 0, len(ents)),
	}
	var (
		err error
		msg *charon.Group
	)
	for _, e := range ents {
		if msg, err = e.Message(); err != nil {
			return nil, err
		}
		resp.Groups = append(resp.Groups, msg)
	}

	return resp, nil
}
