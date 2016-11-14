package model

import "testing"

var (
	userPermissionsTestFixtures = []*UserEntity{
		{
			ID:        1,
			Username:  "user1@example.com",
			FirstName: "first_name_1",
			LastName:  "last_name_1",
			Password:  []byte("0123456789"),
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
			ID:        2,
			Username:  "user2@example.com",
			FirstName: "first_name_2",
			LastName:  "last_name_2",
			Password:  []byte("9876543210"),
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

type userPermissionsFixtures struct {
	got, given UserPermissionsEntity
}

func loadUserPermissionsFixtures(t *testing.T, r UserPermissionsProvider, f []*UserPermissionsEntity) chan userPermissionsFixtures {
	data := make(chan userPermissionsFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.Insert(given)
			if err != nil {
				t.Errorf("user permission cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Log("user permission has been created")
			}

			data <- userPermissionsFixtures{
				got:   *entity,
				given: *given,
			}
		}

		close(data)
	}()

	return data
}
