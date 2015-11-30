package charon

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.True(t, permissions.Contains(expected))
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
		assert.False(t, permissions.Contains(unexpected))
	}
}
