package charon

import "testing"

func TestGetGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req GetGroupRequest
		act actor
	}{
		{
			req: GetGroupRequest{Id: 1},
			act: actor{
				user: &userEntity{ID: 2},
				permissions: Permissions{
					GroupCanRetrieve,
				},
			},
		},
		{
			req: GetGroupRequest{Id: 1},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
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
		req GetGroupRequest
		act actor
	}{
		{
			req: GetGroupRequest{Id: 1},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: GetGroupRequest{Id: 1},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &getGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
