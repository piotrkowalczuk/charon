package charon

import (
	"testing"

	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/qtypes"
)

func TestListUsersHandler_firewall_success(t *testing.T) {
	data := []struct {
		req ListUsersRequest
		act actor
	}{
		{
			req: ListUsersRequest{
				CreatedBy: qtypes.EqualInt64(1),
			},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserCanRetrieveAsOwner,
				},
			},
		},
		{
			req: ListUsersRequest{
				CreatedBy: qtypes.EqualInt64(3),
			},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserCanRetrieveAsStranger,
				},
			},
		},
		{
			req: ListUsersRequest{
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
			req: ListUsersRequest{},
			act: actor{
				user: &userEntity{
					ID:          1,
					IsSuperuser: true,
				},
			},
		},
		{
			req: ListUsersRequest{
				IsStaff:   &ntypes.Bool{Bool: true, Valid: true},
				CreatedBy: qtypes.EqualInt64(1),
			},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: Permissions{
					UserCanRetrieveStaffAsOwner,
				},
			},
		},
		{
			req: ListUsersRequest{
				IsStaff:   &ntypes.Bool{Bool: true, Valid: true},
				CreatedBy: qtypes.EqualInt64(3),
			},
			act: actor{
				user: &userEntity{
					ID: 1,
				},
				permissions: Permissions{
					UserCanRetrieveStaffAsStranger,
				},
			},
		},
		{
			req: ListUsersRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserCanRetrieveAsStranger,
					UserCanRetrieveAsOwner,
					UserCanRetrieveStaffAsStranger,
					UserCanRetrieveStaffAsOwner,
				},
			},
		},
		{
			req: ListUsersRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 1, IsSuperuser: true},
				permissions: Permissions{
					UserCanRetrieveAsStranger,
					UserCanRetrieveAsOwner,
					UserCanRetrieveStaffAsStranger,
					UserCanRetrieveStaffAsOwner,
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
		req ListUsersRequest
		act actor
	}{
		{
			req: ListUsersRequest{},
			act: actor{
				user: &userEntity{},
			},
		},
		{
			req: ListUsersRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: ListUsersRequest{
				IsSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserCanRetrieveAsStranger,
					UserCanRetrieveAsOwner,
					UserCanRetrieveStaffAsStranger,
					UserCanRetrieveStaffAsOwner,
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
