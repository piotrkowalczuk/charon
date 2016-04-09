package charon

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type listGroupsHandler struct {
	*handler
}

func (lgh *listGroupsHandler) handle(ctx context.Context, req *ListGroupsRequest) (*ListGroupsResponse, error) {
	act, err := lgh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = lgh.firewall(req, act); err != nil {
		return nil, err
	}

	ents, err := lgh.repository.group.Find(&groupCriteria{
		offset: req.Offset.Int64Or(0),
		limit:  req.Limit.Int64Or(10),
	})
	if err != nil {
		return nil, err
	}

	return lgh.response(ents)
}

func (lgh *listGroupsHandler) firewall(req *ListGroupsRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(GroupCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charon: list of groups cannot be retrieved, missing permission")
}

func (lgh *listGroupsHandler) response(ents []*groupEntity) (*ListGroupsResponse, error) {
	resp := &ListGroupsResponse{
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
