package charond

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
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for ur := range loadUserFixtures(t, suite.repository.user, userGroupsTestFixtures) {
		for gr := range loadGroupFixtures(t, suite.repository.group, ur.given.Group) {
			add := []*userGroupsEntity{{
				UserID:  ur.got.ID,
				GroupID: gr.got.ID,
			}}
			for range loadUserGroupsFixtures(t, suite.repository.userGroups, add) {
				exists, err := suite.repository.userGroups.Exists(ur.given.ID, gr.given.ID)

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

func TestUserGroupsRepository_Set(t *testing.T) {
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	groups := make([]int64, 0, len(userGroupsTestFixtures))
	for ur := range loadUserFixtures(t, suite.repository.user, userGroupsTestFixtures) {
		for gr := range loadGroupFixtures(t, suite.repository.group, ur.given.Group) {
			groups = append(groups, gr.got.ID)
		}
	}

	i, d, err := suite.repository.userGroups.Set(userGroupsTestFixtures[0].ID, groups)
	if err != nil {
		t.Errorf("user groups cannot be set, unexpected error: %s", err.Error())
	}
	if i != int64(len(groups)) {
		t.Errorf("wrong number of user groups inserted, expected %d but got %d", len(groups), i)
	}
	if d != 0 {
		t.Errorf("wrong number of user groups deleted, expected %d but got %d", 0, d)
	}

	i, d, err = suite.repository.userGroups.Set(userGroupsTestFixtures[0].ID, groups)
	if err != nil {
		t.Errorf("user groups cannot be set, unexpected error: %s", err.Error())
	}
	if i != 0 {
		t.Errorf("wrong number of user groups inserted, expected %d but got %d", 0, i)
	}
	if d != 0 {
		t.Errorf("wrong number of user groups deleted, expected %d but got %d", 0, d)
	}

	i, d, err = suite.repository.userGroups.Set(userGroupsTestFixtures[0].ID, []int64{})
	if err != nil {
		t.Errorf("user groups cannot be set, unexpected error: %s", err.Error())
	}
	if i != 0 {
		t.Errorf("wrong number of user groups inserted, expected %d but got %d", 0, i)
	}
	if d != int64(len(groups)) {
		t.Errorf("wrong number of user groups deleted, expected %d but got %d", len(groups), d)
	}
}

type userGroupsFixtures struct {
	got, given userGroupsEntity
}

func loadUserGroupsFixtures(t *testing.T, r userGroupsProvider, f []*userGroupsEntity) chan userGroupsFixtures {
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
