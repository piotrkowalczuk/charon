package charon

import "testing"

func TestListUserPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req ListUserPermissionsRequest
		act actor
	}{
		{
			req: ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					UserPermissionCanRetrieve,
				},
			},
		},
		{
			req: ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
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
		req ListUserPermissionsRequest
		act actor
	}{
		{
			req: ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{},
				permissions: Permissions{
					PermissionCanRetrieve,
					GroupPermissionCanRetrieve,
				},
			},
		},
	}

	h := &listUserPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
