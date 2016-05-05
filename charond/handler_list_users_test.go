package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/qtypes"
)

func TestListUsersHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.ListUsersRequest
		act actor
	}{
		{
			req: charon.ListUsersRequest{
				CreatedBy: qtypes.EqualInt64(1),
			},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsOwner,
				},
			},
		},
		{
			req: charon.ListUsersRequest{
				CreatedBy: qtypes.EqualInt64(3),
			},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
				},
			},
		},
		{
			req: charon.ListUsersRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
		},
		{
			req: charon.ListUsersRequest{},
			act: actor{
				user: &userEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
		},
		{
			req: charon.ListUsersRequest{
				IsStaff:   &ntypes.Bool{Bool: true, Valid: true},
				CreatedBy: qtypes.EqualInt64(1),
			},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
		},
		{
			req: charon.ListUsersRequest{
				IsStaff:   &ntypes.Bool{Bool: true, Valid: true},
				CreatedBy: qtypes.EqualInt64(3),
			},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: charon.Permissions{
					charon.UserCanRetrieveStaffAsStranger,
				},
			},
		},
		{
			req: charon.ListUsersRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserCanRetrieveAsStranger,
					charon.UserCanRetrieveAsOwner,
					charon.UserCanRetrieveStaffAsStranger,
					charon.UserCanRetrieveStaffAsOwner,
				},
			},
		},
		{
			req: charon.ListUsersRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 1, IsSuperuser: true},
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
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error for %d: %s", i, err.Error())
		}
	}
}

func TestListUsersHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.ListUsersRequest
		act actor
	}{
		{
			req: charon.ListUsersRequest{},
			act: actor{
				user: &userEntity{},
			},
		},
		{
			req: charon.ListUsersRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: charon.ListUsersRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 1},
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
