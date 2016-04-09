package charon

import "testing"

func TestSetUserGroupsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req SetUserGroupsRequest
		act actor
	}{
		{
			req: SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserGroupCanCreate,
					UserGroupCanDelete,
				},
			},
		},
		{
			req: SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &setUserGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestSetUserGroupsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req SetUserGroupsRequest
		act actor
	}{
		{
			req: SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserGroupCanDelete,
				},
			},
		},
		{
			req: SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserGroupCanCreate,
				},
			},
		},
	}

	h := &setUserGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
