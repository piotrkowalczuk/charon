package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestIsGrantedHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.IsGrantedRequest
		act actor
	}{
		{
			req: charon.IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{id: 1},
			},
		},
		{
			req: charon.IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{id: 2},
				permissions: charon.Permissions{
					charon.UserPermissionCanCheckGrantingAsStranger,
				},
			},
		},
		{
			req: charon.IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
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
		req charon.IsGrantedRequest
		act actor
	}{
		{
			req: charon.IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{id: 2},
			},
		},
		{
			req: charon.IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
	}

	h := &isGrantedHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
