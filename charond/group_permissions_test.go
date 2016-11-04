package charond

import (
	"testing"

	"github.com/piotrkowalczuk/ntypes"
)

var (
	groupPermissionsTestFixtures = []*groupEntity{
		{
			id:          1,
			name:        "group_1",
			description: &ntypes.String{String: "first group", Valid: true},
			permissions: []*permissionEntity{
				{
					id:        1,
					subsystem: "subsystem_1",
					module:    "module_1",
					action:    "action_1",
				},
			},
		},
		{
			id:          2,
			name:        "group_2",
			description: &ntypes.String{String: "second group", Valid: true},
			permissions: []*permissionEntity{
				{
					id:        2,
					subsystem: "subsystem_2",
					module:    "module_2",
					action:    "action_2",
				},
			},
		},
	}
)

type groupPermissionsFixtures struct {
	got, given groupPermissionsEntity
}

func loadGroupPermissionsFixtures(t *testing.T, r groupPermissionsProvider, f []*groupPermissionsEntity) chan groupPermissionsFixtures {
	data := make(chan groupPermissionsFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.insert(given)
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
