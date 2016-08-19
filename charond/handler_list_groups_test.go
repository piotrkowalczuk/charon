package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestListGroupsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.ListGroupsRequest
		act actor
	}{
		{
			req: charon.ListGroupsRequest{},
			act: actor{
				user: &userEntity{id: 1},
				permissions: charon.Permissions{
					charon.GroupCanRetrieve,
				},
			},
		},
		{
			req: charon.ListGroupsRequest{},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
			},
		},
	}

	h := &listGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestListGroupsHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.ListGroupsRequest
		act actor
	}{
		{
			req: charon.ListGroupsRequest{},
			act: actor{
				user: &userEntity{id: 1},
			},
		},
		{
			req: charon.ListGroupsRequest{},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
	}

	h := &listGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
