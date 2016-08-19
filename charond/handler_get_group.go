package charond

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type getGroupHandler struct {
	*handler
}

func (ggh *getGroupHandler) handle(ctx context.Context, req *charon.GetGroupRequest) (*charon.GetGroupResponse, error) {
	ggh.loggerWith("group_id", req.Id)

	act, err := ggh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = ggh.firewall(req, act); err != nil {
		return nil, err
	}

	ent, err := ggh.repository.group.findOneByID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "group does not exists")
		}
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	return ggh.response(ent)
}

func (ggh *getGroupHandler) firewall(req *charon.GetGroupRequest, act *actor) error {
	if act.user.isSuperuser {
		return nil
	}
	if act.permissions.Contains(charon.GroupCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "group cannot be retrieved, missing permission")
}

func (ggh *getGroupHandler) response(ent *groupEntity) (*charon.GetGroupResponse, error) {
	msg, err := ent.message()
	if err != nil {
		return nil, err
	}
	return &charon.GetGroupResponse{
		Group: msg,
	}, nil
}
