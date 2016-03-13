// +build unit,!postgres,!e2e

package main

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestSetUserGroupsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.SetUserGroupsRequest
		act actor
	}{
		{
			req: charon.SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserGroupCanCreate,
					charon.UserGroupCanDelete,
				},
			},
		},
		{
			req: charon.SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &setUserGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestSetUserGroupsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.SetUserGroupsRequest
		act actor
	}{
		{
			req: charon.SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: charon.SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: charon.SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserGroupCanDelete,
				},
			},
		},
		{
			req: charon.SetUserGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: charon.Permissions{
					charon.UserGroupCanCreate,
				},
			},
		},
	}

	h := &setUserGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
