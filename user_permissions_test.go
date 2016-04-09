package charon

import "testing"

var (
	userPermissionsTestFixtures = []*userEntity{
		{
			ID:        1,
			Username:  "user1@example.com",
			FirstName: "first_name_1",
			LastName:  "last_name_1",
			Password:  []byte("0123456789"),
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
			ID:        2,
			Username:  "user2@example.com",
			FirstName: "first_name_2",
			LastName:  "last_name_2",
			Password:  []byte("9876543210"),
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

type userPermissionsFixtures struct {
	got, given userPermissionsEntity
}

func loadUserPermissionsFixtures(t *testing.T, r userPermissionsProvider, f []*userPermissionsEntity) chan userPermissionsFixtures {
	data := make(chan userPermissionsFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.Insert(given)
			if err != nil {
				t.Errorf("user permission cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Logf("user permission has been created")
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
