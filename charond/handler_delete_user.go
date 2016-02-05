package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type deleteUserHandler struct {
	*handler
}

func (duh *deleteUserHandler) handle(ctx context.Context, req *charon.DeleteUserRequest) (*charon.DeleteUserResponse, error) {
	duh.loggerWith("user_id", req.Id)

	if req.Id <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "charond: user cannot be deleted, invalid id: %d", req.Id)
	}

	affected, err := duh.repository.user.DeleteByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &charon.DeleteUserResponse{
		Affected: affected,
	}, nil
}
