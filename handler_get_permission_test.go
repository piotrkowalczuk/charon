package charon

import "testing"

func TestGetPermissionHandler_firewall_success(t *testing.T) {
	data := []struct {
		req GetPermissionRequest
		act actor
	}{
		{
			req: GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{ID: 2},
				permissions: Permissions{
					PermissionCanRetrieve,
				},
			},
		},
		{
			req: GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
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
		req GetPermissionRequest
		act actor
	}{
		{
			req: GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{ID: 2},
			},
		},
		{
			req: GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &getPermissionHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
