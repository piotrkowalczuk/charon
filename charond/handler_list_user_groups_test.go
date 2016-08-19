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
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.UserGroupCanRetrieve,
				},
			},
		},
		{
			req: charon.ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
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
				user: &userEntity{id: 1},
			},
		},
		{
			req: charon.ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
		{
			req: charon.ListUserGroupsRequest{},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
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
			t.Error("expected error, got nil")
		}
	}
}
