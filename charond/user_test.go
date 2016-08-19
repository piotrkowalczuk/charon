package charond

import (
	"testing"

	"reflect"

	"github.com/piotrkowalczuk/ntypes"
	"github.com/piotrkowalczuk/qtypes"
)

var (
	userTestFixtures = []*userEntity{
		{
			username:          "johnsnow@gmail.com",
			password:          []byte("secret"),
			firstName:         "John",
			lastName:          "Snow",
			confirmationToken: []byte("1234567890"),
		},
		{
			username:          "1",
			password:          []byte("2"),
			firstName:         "3",
			lastName:          "4",
			confirmationToken: []byte("5"),
		},
	}
)

func TestUserEntity_String(t *testing.T) {
	cases := map[string]userEntity{
		"John Snow": {
			firstName: "John",
			lastName:  "Snow",
		},
		"Snow": {
			lastName: "Snow",
		},
		"John": {
			firstName: "John",
		},
	}

	for expected, c := range cases {
		if c.String() != expected {
			t.Errorf("wrong output, expected %s but got %s", expected, c.String())
		}
	}
}

func TestUserRepository_Create(t *testing.T) {
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for res := range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		if res.got.createdAt.IsZero() {
			t.Errorf("invalid created at field, expected valid time but got %v", res.got.createdAt)
		} else {
			t.Logf("user has been properly created at %v", res.got.createdAt)
		}

		if res.got.username != res.given.username {
			t.Errorf("wrong username, expected %s but got %s", res.given.username, res.got.username)
		}
	}
}

func TestUserRepository_updateByID(t *testing.T) {
	suffix := "_modified"
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for res := range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		user := res.got
		modified, err := suite.repository.user.updateOneByID(user.id, &userPatch{
			firstName:   &ntypes.String{String: user.firstName + suffix, Valid: true},
			isActive:    &ntypes.Bool{Bool: true, Valid: true},
			isConfirmed: &ntypes.Bool{Bool: true, Valid: true},
			isStaff:     &ntypes.Bool{Bool: true, Valid: true},
			isSuperuser: &ntypes.Bool{Bool: true, Valid: true},
			lastName:    &ntypes.String{String: user.lastName + suffix, Valid: true},
			password:    user.password,
			username:    &ntypes.String{String: user.username + suffix, Valid: true},
		})

		if err != nil {
			t.Errorf("user cannot be modified, unexpected error: %s", err.Error())
			continue
		} else {
			t.Logf("user with id %d has been modified", modified.id)
		}

		if assertfTime(t, modified.updatedAt, "invalid updated at field, expected valid time but got %v", modified.updatedAt) {
			t.Logf("user has been properly modified at %v", modified.updatedAt)
		}
		assertf(t, modified.username == user.username+suffix, "wrong username, expected %s but got %s", user.username+suffix, modified.username)
		assertf(t, reflect.DeepEqual(modified.password, user.password), "wrong password, expected %s but got %s", user.password, modified.password)
		assertf(t, modified.firstName == user.firstName+suffix, "wrong first name, expected %s but got %s", user.firstName+suffix, modified.firstName)
		assertf(t, modified.lastName == user.lastName+suffix, "wrong last name, expected %s but got %s", user.lastName+suffix, modified.lastName)
		assert(t, modified.isSuperuser, "user should become a superuser")
		assert(t, modified.isActive, "user should be active")
		assert(t, modified.isConfirmed, "user should be confirmed")
		assert(t, modified.isStaff, "user should become a staff")
	}
}

func TestUserRepository_DeleteByID(t *testing.T) {
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for res := range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		affected, err := suite.repository.user.deleteOneByID(res.got.id)
		if err != nil {
			t.Errorf("user cannot be deleted, unexpected error: %s", err.Error())
			continue
		}

		assert(t, affected == 1, "user was not deleted, no rows affected")

	}
}

func TestUserRepository_findOneByID(t *testing.T) {
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for res := range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		user := res.got
		found, err := suite.repository.user.findOneByID(user.id)

		if err != nil {
			t.Errorf("user cannot be found, unexpected error: %s", err.Error())
			continue
		}

		if assert(t, found != nil, "user was not found, nil object returned") {
			assertf(t, reflect.DeepEqual(res.got, *found), "created and retrieved entity should be equal, but its not\ncreated: %#v\nfounded: %#v", res.got, found)
		}
	}
}

func TestUserRepository_updateLastLoginAt(t *testing.T) {
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for res := range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		affected, err := suite.repository.user.updateLastLoginAt(res.got.id)

		if err != nil {
			t.Errorf("user cannot be updated, unexpected error: %s", err.Error())
			continue
		}

		if assert(t, affected == 1, "user was not updated, no rows affected") {
			entity, err := suite.repository.user.findOneByID(res.got.id)
			if err != nil {
				t.Errorf("user cannot be found, unexpected error: %s", err.Error())
				continue
			}

			assertfTime(t, entity.lastLoginAt, "user last login at property was not properly updated, got %v", entity.lastLoginAt)
		}
	}
}

func TestUserRepository_find(t *testing.T) {
	var (
		err      error
		entities []*userEntity
		all      int64
	)
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		all++
	}

	entities, err = suite.repository.user.find(&userCriteria{limit: all})
	if err != nil {
		t.Errorf("users can not be retrieved, unexpected error: %s", err.Error())
	} else {
		assertf(t, int64(len(entities)) == all, "number of users retrived do not match, expected %d got %d", all, len(entities))
	}

	entities, err = suite.repository.user.find(&userCriteria{
		offset: all,
		limit:  all,
	})
	if err != nil {
		t.Errorf("users can not be retrieved, unexpected error: %s", err.Error())
	} else {
		assertf(t, len(entities) == 0, "number of users retrived do not match, expected %d got %d", 0, len(entities))
	}
	entities, err = suite.repository.user.find(&userCriteria{
		limit:             all,
		username:          qtypes.EqualString("johnsnow@gmail.com"),
		password:          []byte("secret"),
		firstName:         qtypes.EqualString("John"),
		lastName:          qtypes.EqualString("Snow"),
		confirmationToken: []byte("1234567890"),
		isSuperuser:       &ntypes.Bool{Bool: true}, // Is not valid, should not affect results
	})
	if err != nil {
		t.Errorf("users can not be retrieved, unexpected error: %s", err.Error())
	} else {
		assertf(t, len(entities) == 1, "number of users retrived do not match, expected %d got %d", 1, len(entities))
	}
}

func TestUserRepository_IsGranted(t *testing.T) {
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for ur := range loadUserFixtures(t, suite.repository.user, userPermissionsTestFixtures) {
		for pr := range loadPermissionFixtures(t, suite.repository.permission, ur.given.permission) {
			add := []*userPermissionsEntity{{
				userID:              ur.got.id,
				permissionSubsystem: pr.got.subsystem,
				permissionModule:    pr.got.module,
				permissionAction:    pr.got.action,
			}}
			for range loadUserPermissionsFixtures(t, suite.repository.userPermissions, add) {
				exists, err := suite.repository.user.IsGranted(ur.given.id, pr.given.Permission())

				if err != nil {
					t.Errorf("user permission cannot be found, unexpected error: %s", err.Error())
					continue
				}

				if !exists {
					t.Errorf("user permission not found for user %d and permission %d", ur.given.id, pr.given.id)
				} else {
					t.Logf("user permission relationship exists for user %d and permission %d", ur.given.id, pr.given.id)
				}
			}
		}
	}
}

func TestUserRepository_SetPermissions(t *testing.T) {
	t.Skip("not implemented")
}

type userFixtures struct {
	got, given userEntity
}

func loadUserFixtures(t *testing.T, r userProvider, f []*userEntity) chan userFixtures {
	data := make(chan userFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.insert(given)
			if err != nil {
				t.Errorf("user cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Logf("user has been created, got id %d", entity.id)
			}

			data <- userFixtures{
				got:   *entity,
				given: *given,
			}
		}

		close(data)
	}()

	return data
}
