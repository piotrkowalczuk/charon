package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestListUserGroupsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.ListUserGroupsRequest
		act actor
	}{
		{
			req: charon.ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserGroupCanRetrieve,
				},
			},
		},
		{
			req: charon.ListUserGroupsRequest{},
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
		req charon.ListUserGroupsRequest
		act actor
	}{
		{
			req: charon.ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: charon.ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charon.ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
				permissions: charon.Permissions{
					charon.GroupCanRetrieve,
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
