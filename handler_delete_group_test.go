package charon

import "testing"

func TestDeleteGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req DeleteGroupRequest
		act actor
	}{
		{
			req: DeleteGroupRequest{},
			act: actor{
				user: &userEntity{ID: 1},
				permissions: Permissions{
					GroupCanDelete,
				},
			},
		},
		{
			req: DeleteGroupRequest{},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
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
		req DeleteGroupRequest
		act actor
	}{
		{
			req: DeleteGroupRequest{},
			act: actor{
				user: &userEntity{ID: 2},
			},
		},
		{
			req: DeleteGroupRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &deleteGroupHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
