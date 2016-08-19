package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestListUserPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.ListUserPermissionsRequest
		act actor
	}{
		{
			req: charon.ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.UserPermissionCanRetrieve,
				},
			},
		},
		{
			req: charon.ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
			},
		},
	}

	h := &listUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestListUserPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.ListUserPermissionsRequest
		act actor
	}{
		{
			req: charon.ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 1},
			},
		},
		{
			req: charon.ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
		{
			req: charon.ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{},
				permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
					charon.GroupPermissionCanRetrieve,
				},
			},
		},
	}

	h := &listUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
