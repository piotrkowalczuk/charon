// +build unit,!postgres

package main

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
				user: &userEntity{ID: 2},
				permissions: charon.Permissions{
					charon.PermissionCanRetrieve,
				},
			},
		},
		{
			req: charon.GetPermissionRequest{Id: 1},
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
		req charon.GetPermissionRequest
		act actor
	}{
		{
			req: charon.GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: charon.GetPermissionRequest{Id: 1},
			act: actor{
				user: &userEntity{ID: 2},
			},
		},
		{
			req: charon.GetPermissionRequest{Id: 1},
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
