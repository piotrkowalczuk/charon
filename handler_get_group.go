package charon

import (
	"database/sql"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type getGroupHandler struct {
	*handler
}

func (ggh *getGroupHandler) handle(ctx context.Context, req *GetGroupRequest) (*GetGroupResponse, error) {
	ggh.loggerWith("group_id", req.Id)

	act, err := ggh.retrieveActor(ctx)
	if err != nil {
		return nil, err
	}
	if err = ggh.firewall(req, act); err != nil {
		return nil, err
	}

	ent, err := ggh.repository.group.FindOneByID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "charon: group with id %d does not exists", req.Id)
		}
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	return ggh.response(ent)
}

func (ggh *getGroupHandler) firewall(req *GetGroupRequest, act *actor) error {
	if act.user.IsSuperuser {
		return nil
	}
	if act.permissions.Contains(GroupCanRetrieve) {
		return nil
	}

	return grpc.Errorf(codes.PermissionDenied, "charon: group cannot be retrieved, missing permission")
}

func (ggh *getGroupHandler) response(ent *groupEntity) (*GetGroupResponse, error) {
	msg, err := ent.Message()
	if err != nil {
		return nil, err
	}
	return &GetGroupResponse{
		Group: msg,
	}, nil
}
