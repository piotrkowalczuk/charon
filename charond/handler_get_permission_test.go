package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestGetPermissionHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.GetPermissionRequest
		act actor
	}{
		{
			req: charon.GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{id: 2},
				permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
				},
			},
		},
		{
			req: charon.GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
			},
		},
	}

	h := &getPermissionHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestGetPermissionHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.GetPermissionRequest
		act actor
	}{
		{
			req: charon.GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{id: 1},
			},
		},
		{
			req: charon.GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{id: 2},
			},
		},
		{
			req: charon.GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{
					id: 2, isStaff: true,
				},
			},
		},
	}

	h := &getPermissionHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
