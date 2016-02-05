package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
)

type listUsersHandler struct {
	*handler
}

func (luh *listUsersHandler) handle(ctx context.Context, req *charon.ListUsersRequest) (*charon.ListUsersResponse, error) {
	users, err := luh.repository.user.Find(&userCriteria{offset: req.Offset.Int64Or(0), limit: req.Limit.Int64Or(10)})
	if err != nil {
		return nil, err
	}

	luh.loggerWith("results", len(users))

	resp := &charon.ListUsersResponse{
		Users: make([]*charon.User, 0, len(users)),
	}
	for _, u := range users {
		resp.Users = append(resp.Users, u.Message())
	}

	return resp, nil
}
