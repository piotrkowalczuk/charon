package charond

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type deleteGroupHandler struct {
	*handler
}

func (dgh *deleteGroupHandler) handle(ctx context.Context, req *charon.DeleteGroupRequest) (*charon.DeleteGroupResponse, error) {
	dgh.loggerWith("group_id", req.Id)

	act, err := dgh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = dgh.firewall(req, act); err != nil {
		return nil, err
	}

	affected, err := dgh.repository.group.DeleteByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &charon.DeleteGroupResponse{
		Affected: affected,
	}, nil
}

func (dgh *deleteGroupHandler) firewall(req *charon.DeleteGroupRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.GroupCanDelete) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charond: group cannot be removed, missing permission")
}
