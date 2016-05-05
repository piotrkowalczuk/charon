package charond

import (
	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/pqt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type createGroupHandler struct {
	*handler
}

func (cgh *createGroupHandler) handle(ctx context.Context, req *charon.CreateGroupRequest) (*charon.CreateGroupResponse, error) {
	act, err := cgh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = cgh.firewall(req, act); err != nil {
		return nil, err
	}

	ent, err := cgh.repository.group.Create(act.user.ID, req.Name, req.Description)
	if err != nil {
		switch pqt.ErrorConstraint(err) {
		case tableGroupConstraintNameUnique:
			return nil, grpc.Errorf(codes.AlreadyExists, "charond: group with given name already exists")
		default:
			return nil, err
		}
	}

	return cgh.response(ent)
}

func (cgh *createGroupHandler) firewall(req *charon.CreateGroupRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.GroupCanCreate) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charond: group cannot be created, missing permission")
}

func (cgh *createGroupHandler) response(ent *groupEntity) (*charon.CreateGroupResponse, error) {
	msg, err := ent.Message()
	if err != nil {
		return nil, err
	}
	return &charon.CreateGroupResponse{
		Group: msg,
	}, nil
}
