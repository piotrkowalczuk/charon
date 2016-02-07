// +build postgres !unit

package main

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

func TestUserPermissionsRepository_Exists(t *testing.T) {
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	for ur := range loadUserFixtures(t, suite.user, userPermissionsTestFixtures) {
		for pr := range loadPermissionFixtures(t, suite.permission, ur.given.Permission) {
			add := []*userPermissionsEntity{{
				UserID:       ur.got.ID,
				PermissionID: pr.got.ID,
			}}
			for _ = range loadUserPermissionsFixtures(t, suite.userPermissions, add) {
				exists, err := suite.userPermissions.Exists(ur.given.ID, pr.given.Permission())

				if err != nil {
					t.Errorf("user permission cannot be found, unexpected error: %s", err.Error())
					continue
				}

				if !exists {
					t.Errorf("user permission not found for user %d and permission %d", ur.given.ID, pr.given.ID)
				} else {
					t.Logf("user permission relationship exists for user %d and permission %d", ur.given.ID, pr.given.ID)
				}
			}
		}
	}
}

type userPermissionsFixtures struct {
	got, given userPermissionsEntity
}

func loadUserPermissionsFixtures(t *testing.T, r UserPermissionsRepository, f []*userPermissionsEntity) chan userPermissionsFixtures {
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
