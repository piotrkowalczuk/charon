package main

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
				user: &userEntity{ID: 1},
			},
		},
		{
			req: charon.IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{ID: 2},
				permissions: charon.Permissions{
					charon.UserPermissionCanCheckGrantingAsStranger,
				},
			},
		},
		{
			req: charon.IsGrantedRequest{UserId: 1},
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
		req charon.IsGrantedRequest
		act actor
	}{
		{
			req: charon.IsGrantedRequest{UserId: 1},
			act: actor{
				user: &userEntity{ID: 2},
			},
		},
		{
			req: charon.IsGrantedRequest{UserId: 1},
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
