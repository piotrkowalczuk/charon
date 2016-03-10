// +build unit !postgres

package main

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
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserPermissionCanRetrieve,
				},
			},
		},
		{
			req: charon.ListUserPermissionsRequest{},
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
		req charon.ListUserPermissionsRequest
		act actor
	}{
		{
			req: charon.ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: charon.ListUserPermissionsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
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
			t.Errorf("expected error, got nil")
		}
	}
}
