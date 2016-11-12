package charond

import (
	"database/sql"

	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type modifyGroupHandler struct {
	*handler
}

// TODO: missing firewall
func (mgh *modifyGroupHandler) Modify(ctx context.Context, req *charonrpc.ModifyGroupRequest) (*charonrpc.ModifyGroupResponse, error) {
	mgh.loggerWith("group_id", req.Id)

	actor, err := mgh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}

	mgh.loggerWith("user_id", actor.user.ID)

	group, err := mgh.repository.group.UpdateOneByID(req.Id, actor.user.ID, req.Name, req.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "group does not exists")
		}
		return nil, err
	}

	return mgh.response(group)
}

func (mgh *modifyGroupHandler) response(g *model.GroupEntity) (*charonrpc.ModifyGroupResponse, error) {
	msg, err := g.Message()
	if err != nil {
		return nil, err
	}
	return &charonrpc.ModifyGroupResponse{
		Group: msg,
	}, nil
}
