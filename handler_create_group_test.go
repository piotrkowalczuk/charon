package charon

import "testing"

func TestCreateGroupHandler_firewall_success(t *testing.T) {
	data := []struct {
		req CreateGroupRequest
		act actor
	}{
		{
			req: CreateGroupRequest{},
			act: actor{
				user: &userEntity{ID: 2},
				permissions: Permissions{
					GroupCanCreate,
				},
			},
		},
		{
			req: CreateGroupRequest{},
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
		req CreateGroupRequest
		act actor
	}{
		{
			req: CreateGroupRequest{},
			act: actor{
				user: &userEntity{ID: 2},
			},
		},
		{
			req: CreateGroupRequest{},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
		{
			req: CreateGroupRequest{},
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
