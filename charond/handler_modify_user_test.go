package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
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
			hint:   "superuser should be able to promote an user",
			req:    charonrpc.ModifyUserRequest{Id: 2, IsSuperuser: &ntypes.Bool{Bool: true, Valid: true}},
			entity: model.UserEntity{ID: 2},
			user:   model.UserEntity{ID: 1, IsSuperuser: true},
		},
		{
			hint:   "superuser should be able to promote a staff user",
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
			hint:   "if user has permission to modify an user as a stranger, he should be able to do that",
			req:    charonrpc.ModifyUserRequest{Id: 2},
			entity: model.UserEntity{ID: 2},
			user:   model.UserEntity{ID: 1},
			perm:   charon.Permissions{charon.UserCanModifyAsStranger},
		},
		{
			hint:   "if user has permission to modify an user as an owner, he should be able to do that",
			req:    charonrpc.ModifyUserRequest{Id: 1},
			entity: model.UserEntity{ID: 1},
			user:   model.UserEntity{ID: 1},
			perm:   charon.Permissions{charon.UserCanModifyAsOwner},
		},
	}

	handler := &modifyUserHandler{}
	for _, args := range success {
		msg, ok := handler.firewall(&args.req, &args.entity, &actor{
			user:        &args.user,
			permissions: args.perm,
		})
		if !ok {
			t.Errorf(args.hint+", got: %s", msg)
		} else {
			t.Log(args.hint)
		}
	}
}
