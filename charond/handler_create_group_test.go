// +build unit,!postgres,!e2e

package main

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestCreateGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.CreateGroupRequest
		act actor
	}{
		{
			req: charon.CreateGroupRequest{},
			act: actor{
				user: &userEntity{ID: 2},
				permissions: charon.Permissions{
					charon.GroupCanCreate,
				},
			},
		},
		{
			req: charon.CreateGroupRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &createGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestCreateGroupHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.CreateGroupRequest
		act actor
	}{
		{
			req: charon.CreateGroupRequest{},
			act: actor{
				user: &userEntity{ID: 2},
			},
		},
		{
			req: charon.CreateGroupRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charon.CreateGroupRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
	}

	h := &createGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
