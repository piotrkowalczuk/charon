package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charonrpc"
	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/qtypes"
)

func TestListUsersHandler_firewall_success(t *testing.T) {
	cases := map[string]struct {
		req charonrpc.ListUsersRequest
		act actor
	}{
		"as-owner": {

			req: charonrpc.ListUsersRequest{
				CreatedBy: qtypes.EqualInt64(1),
			},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsOwner,
				},
			},
		},
		"as-stranger": {
			req: charonrpc.ListUsersRequest{
				CreatedBy: qtypes.EqualInt64(3),
			},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
				},
			},
		},
		"as-superuser-search-for-superusers": {
			req: charonrpc.ListUsersRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &model.UserEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
		},
		"as-superuser": {
			req: charonrpc.ListUsersRequest{},
			act: actor{
				user: &model.UserEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
		},
		"as-owner-search-for-staff": {
			req: charonrpc.ListUsersRequest{
				IsStaff:   &ntypes.Bool{Bool: true, Valid: true},
				CreatedBy: qtypes.EqualInt64(1),
			},
			act: actor{
				user: &model.UserEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
		},
		"as-stranger-search-for-staff": {
			req: charonrpc.ListUsersRequest{
				IsStaff:   &ntypes.Bool{Bool: true, Valid: true},
				CreatedBy: qtypes.EqualInt64(3),
			},
			act: actor{
				user: &model.UserEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsStranger,
				},
			},
		},
		"all-permissions": {
			req: charonrpc.ListUsersRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
		},
		"as-superuser-with-all-permissions": {
			req: charonrpc.ListUsersRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &model.UserEntity{ID: 1, IsSuperuser: true},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
		},
	}

	h := &listUsersHandler{}
	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			if err := h.firewall(&c.req, &c.act); err != nil {
				t.Errorf("unexpected error for %d: %s", i, err.Error())
			}
		})
	}
}

func TestListUsersHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charonrpc.ListUsersRequest
		act actor
	}{
		{
			req: charonrpc.ListUsersRequest{},
			act: actor{
				user: &model.UserEntity{},
			},
		},
		{
			req: charonrpc.ListUsersRequest{},
			act: actor{
				user: &model.UserEntity{ID: 1},
			},
		},
		{
			req: charonrpc.ListUsersRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &model.UserEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
		},
	}

	h := &listUsersHandler{}
	for i, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error for %d, got nil", i)
		}
	}
}
