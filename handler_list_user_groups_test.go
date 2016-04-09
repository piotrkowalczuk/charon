package charon

import "testing"

func TestListUserGroupsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req ListUserGroupsRequest
		act actor
	}{
		{
			req: ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserGroupCanRetrieve,
				},
			},
		},
		{
			req: ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &listUserGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestListUserGroupsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req ListUserGroupsRequest
		act actor
	}{
		{
			req: ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
				permissions: Permissions{
					GroupCanRetrieve,
				},
			},
		},
	}

	h := &listUserGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
