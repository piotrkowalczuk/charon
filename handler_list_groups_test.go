package charon

import "testing"

func TestListGroupsHandler_firewall_success(t *testing.T) {
	data := []struct {
		req ListGroupsRequest
		act actor
	}{
		{
			req: ListGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					GroupCanRetrieve,
				},
			},
		},
		{
			req: ListGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
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
		req ListGroupsRequest
		act actor
	}{
		{
			req: ListGroupsRequest{},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: ListGroupsRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &listGroupsHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
