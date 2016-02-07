// +build postgres !unit

package main

import "testing"

var (
	userGroupsTestFixtures = []*userEntity{
		{
			ID:        1,
			Username:  "user1@example.com",
			FirstName: "first_name_1",
			LastName:  "last_name_1",
			Password:  []byte("0123456789"),
			Group: []*groupEntity{
				{
					ID:   1,
					Name: "group_1",
				},
			},
		},
		{
			ID:        2,
			Username:  "user2@example.com",
			FirstName: "first_name_2",
			LastName:  "last_name_2",
			Password:  []byte("9876543210"),
			Group: []*groupEntity{
				{
					ID:   2,
					Name: "group_2",
				},
			},
		},
	}
)

func TestUserGroupsRepository_Exists(t *testing.T) {
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	for ur := range loadUserFixtures(t, suite.user, userGroupsTestFixtures) {
		for gr := range loadGroupFixtures(t, suite.group, ur.given.Group) {
			add := []*userGroupsEntity{{
				UserID:  ur.got.ID,
				GroupID: gr.got.ID,
			}}
			for _ = range loadUserGroupsFixtures(t, suite.userGroups, add) {
				exists, err := suite.userGroups.Exists(ur.given.ID, gr.given.ID)

				if err != nil {
					t.Errorf("user group cannot be found, unexpected error: %s", err.Error())
					continue
				}

				if !exists {
					t.Errorf("user group not found for user %d and group %d", ur.given.ID, gr.given.ID)
				} else {
					t.Logf("user group relationship exists for user %d and group %d", ur.given.ID, gr.given.ID)
				}
			}
		}
	}
}

type userGroupsFixtures struct {
	got, given userGroupsEntity
}

func loadUserGroupsFixtures(t *testing.T, r UserGroupsRepository, f []*userGroupsEntity) chan userGroupsFixtures {
	data := make(chan userGroupsFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.Insert(given)
			if err != nil {
				t.Errorf("user group cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Logf("user group has been created")
			}

			data <- userGroupsFixtures{
				got:   *entity,
				given: *given,
			}
		}

		close(data)
	}()

	return data
}
