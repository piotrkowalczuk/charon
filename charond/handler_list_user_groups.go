package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listUserGroupsHandler struct {
	*handler
}

// TODO: missing firewall
func (lugh *listUserGroupsHandler) ListGroups(ctx context.Context, req *charonrpc.ListUserGroupsRequest) (*charonrpc.ListUserGroupsResponse, error) {
	ents, err := lugh.repository.group.FindByUserID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return lugh.response(ents)
}

func (lugh *listUserGroupsHandler) firewall(req *charonrpc.ListUserGroupsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.UserGroupCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "list of user groups cannot be retrieved, missing permission")
}

func (lugh *listUserGroupsHandler) response(ents []*model.GroupEntity) (*charonrpc.ListUserGroupsResponse, error) {
	resp := &charonrpc.ListUserGroupsResponse{
		Groups: make([]*charonrpc.Group, 0, len(ents)),
	}
	var (
		err error
		msg *charonrpc.Group
	)
	for _, e := range ents {
		if msg, err = e.Message(); err != nil {
			return nil, err
		}
		resp.Groups = append(resp.Groups, msg)
	}

	return resp, nil
}
