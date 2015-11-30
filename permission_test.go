package charon

import (
	"testing"
)

func TestPermissions_Contains(t *testing.T) {
	positive := map[Permission]Permissions{
		UserCanCreate: {
			UserCanCreate,
			UserCanDeleteAsOwner,
		},
		UserPermissionCanCreate: {
			UserPermissionCanCreate,
		},
	}

	for expected, permissions := range positive {
		if permissions.Contains(unexpected) {
			t.Errorf("expected permission (%s), is not present", expected)
		}
	}

	negative := map[Permission]Permissions{
		PermissionCanCreate: {
			UserCanCreate,
			UserCanDeleteAsOwner,
		},
		UserCanDeleteAsStranger: {
			UserPermissionCanCreate,
		},
		UserCanCreate: {},
	}

	for unexpected, permissions := range negative {
		if permissions.Contains(unexpected) {
			t.Errorf("unexpected permission (%s), should not be present", unexpected)
		}
	}
}
