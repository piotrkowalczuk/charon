package charon

import "testing"

func TestListPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req ListPermissionsRequest
		act actor
	}{
		{
			req: ListPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					PermissionCanRetrieve,
				},
			},
		},
		{
			req: ListPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &listPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestListPermissionsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req ListPermissionsRequest
		act actor
	}{
		{
			req: ListPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: ListPermissionsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &listPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
