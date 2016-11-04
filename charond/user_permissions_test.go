package charond

import "testing"

var (
	userPermissionsTestFixtures = []*userEntity{
		{
			id:        1,
			username:  "user1@example.com",
			firstName: "first_name_1",
			lastName:  "last_name_1",
			password:  []byte("0123456789"),
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
			id:        2,
			username:  "user2@example.com",
			firstName: "first_name_2",
			lastName:  "last_name_2",
			password:  []byte("9876543210"),
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

type userPermissionsFixtures struct {
	got, given userPermissionsEntity
}

func loadUserPermissionsFixtures(t *testing.T, r userPermissionsProvider, f []*userPermissionsEntity) chan userPermissionsFixtures {
	data := make(chan userPermissionsFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.insert(given)
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
