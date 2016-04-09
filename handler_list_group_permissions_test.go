package charon

import "testing"

func TestListGroupPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req ListGroupPermissionsRequest
		act actor
	}{
		{
			req: ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					GroupPermissionCanRetrieve,
				},
			},
		},
		{
			req: ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
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
		req ListGroupPermissionsRequest
		act actor
	}{
		{
			req: ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{},
				permissions: Permissions{
					PermissionCanRetrieve,
					UserPermissionCanRetrieve,
				},
			},
		},
	}

	h := &listGroupPermissionsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
