package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
)

type getUserHandler struct {
	*handler
}

func (guh *getUserHandler) handle(ctx context.Context, req *charon.GetUserRequest) (*charon.GetUserResponse, error) {
	guh.loggerWith("user_id", req.Id)

	user, err := guh.repository.user.FindOneByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &charon.GetUserResponse{
		User: user.Message(),
	}, nil
}
