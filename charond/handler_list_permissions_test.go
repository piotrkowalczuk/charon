package main

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestListPermissionsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.ListPermissionsRequest
		act actor
	}{
		{
			req: charon.ListPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
				},
			},
		},
		{
			req: charon.ListPermissionsRequest{},
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
		req charon.ListPermissionsRequest
		act actor
	}{
		{
			req: charon.ListPermissionsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: charon.ListPermissionsRequest{},
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
