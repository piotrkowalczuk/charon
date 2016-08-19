package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestDeleteGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.DeleteGroupRequest
		act actor
	}{
		{
			req: charon.DeleteGroupRequest{},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.GroupCanDelete,
				},
			},
		},
		{
			req: charon.DeleteGroupRequest{},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
			},
		},
	}

	h := &deleteGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestDeleteGroupHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.DeleteGroupRequest
		act actor
	}{
		{
			req: charon.DeleteGroupRequest{},
			act: actor{
				user: &userEntity{id: 2},
			},
		},
		{
			req: charon.DeleteGroupRequest{},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
	}

	h := &deleteGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
