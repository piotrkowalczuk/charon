package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestGetGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.GetGroupRequest
		act actor
	}{
		{
			req: charon.GetGroupRequest{Id: 1},
			act: actor{
				user: &userEntity{id: 2},
				permissions: charon.Permissions{
					charon.GroupCanRetrieve,
				},
			},
		},
		{
			req: charon.GetGroupRequest{Id: 1},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
			},
		},
	}

	h := &getGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestGetGroupHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.GetGroupRequest
		act actor
	}{
		{
			req: charon.GetGroupRequest{Id: 1},
			act: actor{
				user: &userEntity{id: 1},
			},
		},
		{
			req: charon.GetGroupRequest{Id: 1},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
	}

	h := &getGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
