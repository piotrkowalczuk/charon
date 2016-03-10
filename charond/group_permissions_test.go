// +build postgres !unit

package main

import "testing"

var (
	groupPermissionsTestFixtures = []*groupEntity{
		{
			ID:   1,
			Name: "group_1",
			Permission: []*permissionEntity{
				{
					ID:        1,
					Subsystem: "subsystem_1",
					Module:    "module_1",
					Action:    "action_1",
				},
			},
		},
		{
			ID:   2,
			Name: "group_2",
			Permission: []*permissionEntity{
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
	got, given groupPermissionsEntity
}

func loadGroupPermissionsFixtures(t *testing.T, r GroupPermissionsRepository, f []*groupPermissionsEntity) chan groupPermissionsFixtures {
	data := make(chan groupPermissionsFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.Insert(given)
			if err != nil {
				t.Errorf("group permission cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Logf("group permission has been created")
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
