package charon

import "testing"

func TestSetUserPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req SetUserPermissionsRequest
		act actor
	}{
		{
			req: SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserPermissionCanDelete,
					UserPermissionCanCreate,
				},
			},
		},
		{
			req: SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
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
		req SetUserPermissionsRequest
		act actor
	}{
		{
			req: SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserPermissionCanDelete,
				},
			},
		},
		{
			req: SetUserPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserPermissionCanCreate,
				},
			},
		},
	}

	h := &setUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
