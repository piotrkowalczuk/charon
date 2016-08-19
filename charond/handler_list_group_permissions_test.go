package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestListGroupPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.ListGroupPermissionsRequest
		act actor
	}{
		{
			req: charon.ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.GroupPermissionCanRetrieve,
				},
			},
		},
		{
			req: charon.ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
			},
		},
	}

	h := &listGroupPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestListGroupPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.ListGroupPermissionsRequest
		act actor
	}{
		{
			req: charon.ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 1},
			},
		},
		{
			req: charon.ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
		{
			req: charon.ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{},
				permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
					charon.UserPermissionCanRetrieve,
				},
			},
		},
	}

	h := &listGroupPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
