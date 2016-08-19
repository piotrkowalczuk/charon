package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestSetUserPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.SetUserPermissionsRequest
		act actor
	}{
		{
			req: charon.SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.UserPermissionCanDelete,
					charon.UserPermissionCanCreate,
				},
			},
		},
		{
			req: charon.SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
			},
		},
	}

	h := &setUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestSetUserPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.SetUserPermissionsRequest
		act actor
	}{
		{
			req: charon.SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 1},
			},
		},
		{
			req: charon.SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
		{
			req: charon.SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.UserPermissionCanDelete,
				},
			},
		},
		{
			req: charon.SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.UserPermissionCanCreate,
				},
			},
		},
	}

	h := &setUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
