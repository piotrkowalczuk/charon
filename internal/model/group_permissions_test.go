package model

import (
	"context"
	"testing"

	"github.com/piotrkowalczuk/ntypes"
)

var (
	groupPermissionsTestFixtures = []*GroupEntity{
		{
			ID:          1,
			Name:        "group_1",
			Description: ntypes.String{String: "first group", Valid: true},
			Permissions: []*PermissionEntity{
				{
					ID:        1,
					Subsystem: "subsystem_1",
					Module:    "module_1",
					Action:    "action_1",
				},
			},
		},
		{
			ID:          2,
			Name:        "group_2",
			Description: ntypes.String{String: "second group", Valid: true},
			Permissions: []*PermissionEntity{
				{
					ID:        2,
					Subsystem: "subsystem_2",
					Module:    "module_2",
					Action:    "action_2",
				},
			},
		},
	}
)

type groupPermissionsFixtures struct {
	got, given GroupPermissionsEntity
}

func loadGroupPermissionsFixtures(t *testing.T, r GroupPermissionsProvider, f []*GroupPermissionsEntity) chan groupPermissionsFixtures {
	data := make(chan groupPermissionsFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.Insert(context.TODO(), given)
			if err != nil {
				t.Errorf("group permission cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Log("group permission has been created")
			}

			data <- groupPermissionsFixtures{
				got:   *entity,
				given: *given,
			}
		}

		close(data)
	}()

	return data
}
