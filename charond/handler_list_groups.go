package charond

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listGroupsHandler struct {
	*handler
}

func (lgh *listGroupsHandler) handle(ctx context.Context, req *charon.ListGroupsRequest) (*charon.ListGroupsResponse, error) {
	act, err := lgh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = lgh.firewall(req, act); err != nil {
		return nil, err
	}

	ents, err := lgh.repository.group.find(&groupCriteria{
		offset: req.Offset.Int64Or(0),
		limit:  req.Limit.Int64Or(10),
	})
	if err != nil {
		return nil, err
	}

	return lgh.response(ents)
}

func (lgh *listGroupsHandler) firewall(req *charon.ListGroupsRequest, act *actor) error {
	if act.user.isSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.GroupCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "list of groups cannot be retrieved, missing permission")
}

func (lgh *listGroupsHandler) response(ents []*groupEntity) (*charon.ListGroupsResponse, error) {
	resp := &charon.ListGroupsResponse{
		Groups: make([]*charon.Group, 0, len(ents)),
	}
	var (
		err error
		msg *charon.Group
	)
	for _, e := range ents {
		if msg, err = e.message(); err != nil {
			return nil, err
		}
		resp.Groups = append(resp.Groups, msg)
	}

	return resp, nil
}
