package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/session"
	"github.com/piotrkowalczuk/ntypes"
)

func TestModifyUserHandler_Firewall(t *testing.T) {
	success := []struct {
		hint   string
		req    charonrpc.ModifyUserRequest
		entity model.UserEntity
		user   model.UserEntity
		perm   charon.Permissions
	}{

		{
			hint:   "superuser should be able to degrade itself",
			req:    charonrpc.ModifyUserRequest{Id: 1, IsSuperuser: &ntypes.Bool{Bool: false, Valid: true}},
			entity: model.UserEntity{ID: 1, IsSuperuser: true},
			user:   model.UserEntity{ID: 1, IsSuperuser: true},
		},
		{
			hint:   "superuser should be able to promote an User",
			req:    charonrpc.ModifyUserRequest{Id: 2, IsSuperuser: &ntypes.Bool{Bool: true, Valid: true}},
			entity: model.UserEntity{ID: 2},
			user:   model.UserEntity{ID: 1, IsSuperuser: true},
		},
		{
			hint:   "superuser should be able to promote a staff User",
			req:    charonrpc.ModifyUserRequest{Id: 2, IsSuperuser: &ntypes.Bool{Bool: true, Valid: true}},
			entity: model.UserEntity{ID: 2},
			user:   model.UserEntity{ID: 1, IsSuperuser: true},
		},
		{
			hint:   "superuser should be able to degrade another superuser",
			req:    charonrpc.ModifyUserRequest{Id: 2, IsSuperuser: &ntypes.Bool{Bool: false, Valid: true}},
			entity: model.UserEntity{ID: 2, IsSuperuser: true},
			user:   model.UserEntity{ID: 1, IsSuperuser: true},
		},
		{
			hint:   "if User has permission to modify an User as a stranger, he should be able to do that",
			req:    charonrpc.ModifyUserRequest{Id: 2},
			entity: model.UserEntity{ID: 2},
			user:   model.UserEntity{ID: 1},
			perm:   charon.Permissions{charon.UserCanModifyAsStranger},
		},
		{
			hint:   "if User has permission to modify an User as an owner, he should be able to do that",
			req:    charonrpc.ModifyUserRequest{Id: 1},
			entity: model.UserEntity{ID: 1},
			user:   model.UserEntity{ID: 1},
			perm:   charon.Permissions{charon.UserCanModifyAsOwner},
		},
	}

	handler := &modifyUserHandler{}
	for _, args := range success {
		msg, ok := handler.firewall(&args.req, &args.entity, &session.Actor{
			User:        &args.user,
			Permissions: args.perm,
		})
		if !ok {
			t.Errorf(args.hint+", got: %s", msg)
		} else {
			t.Log(args.hint)
		}
	}
}
