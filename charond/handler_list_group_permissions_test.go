// +build unit !postgres

package main

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
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.GroupPermissionCanRetrieve,
				},
			},
		},
		{
			req: charon.ListGroupPermissionsRequest{},
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
		req charon.ListGroupPermissionsRequest
		act actor
	}{
		{
			req: charon.ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: charon.ListGroupPermissionsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
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
			t.Errorf("expected error, got nil")
		}
	}
}
