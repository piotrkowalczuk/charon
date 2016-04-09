package charon

import "testing"

func TestIsGrantedHandler_firewall_success(t *testing.T) {
	data := []struct {
		req IsGrantedRequest
		act actor
	}{
		{
			req: IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{ID: 1},
			},
		},
		{
			req: IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{ID: 2},
				permissions: Permissions{
					UserPermissionCanCheckGrantingAsStranger,
				},
			},
		},
		{
			req: IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{ID: 2, IsSuperuser: true},
			},
		},
	}

	h := &isGrantedHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestIsGrantedHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req IsGrantedRequest
		act actor
	}{
		{
			req: IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{ID: 2},
			},
		},
		{
			req: IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{
					ID:      2,
					IsStaff: true,
				},
			},
		},
	}

	h := &isGrantedHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
