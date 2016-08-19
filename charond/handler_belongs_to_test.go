package charond

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
)

func TestBelongsToHandler_firewall_success(t *testing.T) {
	data := []struct {
		req charon.BelongsToRequest
		act actor
	}{
		{
			req: charon.BelongsToRequest{UserId: 1},
			act: actor{
				user: &userEntity{id: 2},
				permissions: charon.Permissions{
					charon.UserGroupCanCheckBelongingAsStranger,
				},
			},
		},
		{
			req: charon.BelongsToRequest{UserId: 1},
			act: actor{
				user: &userEntity{id: 2, isSuperuser: true},
			},
		},
		{
			req: charon.BelongsToRequest{UserId: 1},
			act: actor{
				user: &userEntity{id: 1},
			},
		},
	}

	h := &belongsToHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestBelongsToHandler_firewall_failure(t *testing.T) {
	data := []struct {
		req charon.BelongsToRequest
		act actor
	}{
		{
			req: charon.BelongsToRequest{UserId: 1},
			act: actor{
				user: &userEntity{id: 2},
			},
		},
		{
			req: charon.BelongsToRequest{UserId: 1},
			act: actor{
				user: &userEntity{
					id:      2,
					isStaff: true,
				},
			},
		},
	}

	h := &belongsToHandler{}
	for _, d := range data {
		if err := h.firewall(&d.req, &d.act); err == nil {
			t.Error("expected error, got nil")
		}
	}
}
