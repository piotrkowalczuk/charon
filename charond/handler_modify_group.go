package charond

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type modifyGroupHandler struct {
	*handler
}

// TODO: missing firewall
func (mgh *modifyGroupHandler) handle(ctx context.Context, req *charon.ModifyGroupRequest) (*charon.ModifyGroupResponse, error) {
	mgh.loggerWith("group_id", req.Id)

	actor, err := mgh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	mgh.loggerWith("user_id", actor.user.id)

	group, err := mgh.repository.group.updateOneByID(req.Id, actor.user.id, req.Name, req.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "group does not exists")
		}
		return nil, err
	}

	return mgh.response(group)
}

func (mgh *modifyGroupHandler) response(g *groupEntity) (*charon.ModifyGroupResponse, error) {
	msg, err := g.message()
	if err != nil {
		return nil, err
	}
	return &charon.ModifyGroupResponse{
		Group: msg,
	}, nil
}
