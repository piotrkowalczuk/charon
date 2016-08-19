package charond

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
				user: &userEntity{id: 2},
				permissions: charon.Permissions{
					charon.GroupCanCreate,
				},
			},
		},
		{
			req: charon.CreateGroupRequest{},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
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
				user: &userEntity{id: 2},
			},
		},
		{
			req: charon.CreateGroupRequest{},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
		{
			req: charon.CreateGroupRequest{},
			act: actor{
				user: &userEntity{id: 1},
			},
		},
	}

	h := &createGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
