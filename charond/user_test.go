// +build postgres !unit

package main

import (
	"testing"

	"reflect"

	"github.com/piotrkowalczuk/nilt"
)

var (
	userTestFixtures = []*userEntity{
		{
			Username:          "johnsnow@gmail.com",
			Password:          []byte("secret"),
			FirstName:         "John",
			LastName:          "Snow",
			ConfirmationToken: []byte("1234567890"),
		},
		{
			Username:          "1",
			Password:          []byte("2"),
			FirstName:         "3",
			LastName:          "4",
			ConfirmationToken: []byte("5"),
		},
	}
)

func TestUserRepository_Create(t *testing.T) {
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	for res := range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		if res.got.CreatedAt.IsZero() {
			t.Errorf("invalid created at field, expected valid time but got %v", res.got.CreatedAt)
		} else {
			t.Logf("user has been properly created at %v", res.got.CreatedAt)
		}

		if res.got.Username != res.given.Username {
			t.Errorf("wrong username, expected %s but got %s", res.given.Username, res.got.Username)
		}
	}
}

func TestUserRepository_UpdateByID(t *testing.T) {
	suffix := "_modified"
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	for res := range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		user := res.got
		modified, err := suite.repository.user.UpdateByID(
			user.ID,
			nil,
			nil,
			nilt.Int64{},
			nilt.String{String: user.FirstName + suffix, Valid: true},
			nilt.Bool{Bool: true, Valid: true},
			nilt.Bool{Bool: true, Valid: true},
			nilt.Bool{Bool: true, Valid: true},
			nilt.Bool{Bool: true, Valid: true},
			nil,
			nilt.String{String: user.LastName + suffix, Valid: true},
			user.Password,
			nil,
			nilt.Int64{},
			nilt.String{String: user.Username + suffix, Valid: true},
		)

		if err != nil {
			t.Errorf("user cannot be modified, unexpected error: %s", err.Error())
			continue
		} else {
			t.Logf("user with id %d has been modified", modified.ID)
		}

		if assertfTime(t, modified.UpdatedAt, "invalid updated at field, expected valid time but got %v", modified.UpdatedAt) {
			t.Logf("user has been properly modified at %v", modified.UpdatedAt)
		}
		assertf(t, modified.Username == user.Username+suffix, "wrong username, expected %s but got %s", user.Username+suffix, modified.Username)
		assertf(t, reflect.DeepEqual(modified.Password, user.Password), "wrong password, expected %s but got %s", user.Password, modified.Password)
		assertf(t, modified.FirstName == user.FirstName+suffix, "wrong first name, expected %s but got %s", user.FirstName+suffix, modified.FirstName)
		assertf(t, modified.LastName == user.LastName+suffix, "wrong last name, expected %s but got %s", user.LastName+suffix, modified.LastName)
		assert(t, modified.IsSuperuser, "user should become a superuser")
		assert(t, modified.IsActive, "user should be active")
		assert(t, modified.IsConfirmed, "user should be confirmed")
		assert(t, modified.IsStaff, "user should become a staff")
	}
}

func TestUserRepository_DeleteByID(t *testing.T) {
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	for res := range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		affected, err := suite.repository.user.DeleteByID(res.got.ID)

		if err != nil {
			t.Errorf("user cannot be deleted, unexpected error: %s", err.Error())
			continue
		}

		assert(t, affected == 1, "user was not deleted, no rows affected")

	}
}

func TestUserRepository_FindOneByID(t *testing.T) {
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	for res := range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		user := res.got
		found, err := suite.repository.user.FindOneByID(user.ID)

		if err != nil {
			t.Errorf("user cannot be found, unexpected error: %s", err.Error())
			continue
		}

		if assert(t, found != nil, "user was not found, nil object returned") {
			assertf(t, reflect.DeepEqual(res.got, *found), "created and retrieved entity should be equal, but its not\ncreated: %#v\nfounded: %#v", res.got, found)
		}
	}
}

func TestUserRepository_UpdateLastLoginAt(t *testing.T) {
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	for res := range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		affected, err := suite.repository.user.UpdateLastLoginAt(res.got.ID)

		if err != nil {
			t.Errorf("user cannot be updated, unexpected error: %s", err.Error())
			continue
		}

		if assert(t, affected == 1, "user was not updated, no rows affected") {
			entity, err := suite.repository.user.FindOneByID(res.got.ID)
			if err != nil {
				t.Errorf("user cannot be found, unexpected error: %s", err.Error())
				continue
			}

			assertfTime(t, entity.LastLoginAt, "user last login at property was not properly updated, got %v", entity.LastLoginAt)
		}
	}
}

func TestUserRepository_Find(t *testing.T) {
	var (
		err      error
		entities []*userEntity
		all      int64
	)
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	for _ = range loadUserFixtures(t, suite.repository.user, userTestFixtures) {
		all++
	}

	entities, err = suite.repository.user.Find(&userCriteria{limit: all})
	if err != nil {
		t.Errorf("users can not be retrieved, unexpected error: %s", err.Error())
	} else {
		assertf(t, int64(len(entities)) == all, "number of users retrived do not match, expected %d got %d", all, len(entities))
	}

	entities, err = suite.repository.user.Find(&userCriteria{
		offset: all,
		limit:  all,
	})
	if err != nil {
		t.Errorf("users can not be retrieved, unexpected error: %s", err.Error())
	} else {
		assertf(t, len(entities) == 0, "number of users retrived do not match, expected %d got %d", 0, len(entities))
	}
	entities, err = suite.repository.user.Find(&userCriteria{
		limit:             all,
		username:          nilt.String{String: "johnsnow@gmail.com", Valid: true},
		password:          []byte("secret"),
		firstName:         nilt.String{String: "John", Valid: true},
		lastName:          nilt.String{String: "Snow", Valid: true},
		confirmationToken: []byte("1234567890"),
		isSuperuser:       nilt.Bool{Bool: true}, // Is not valid, should not affect results
	})
	if err != nil {
		t.Errorf("users can not be retrieved, unexpected error: %s", err.Error())
	} else {
		assertf(t, len(entities) == 1, "number of users retrived do not match, expected %d got %d", 1, len(entities))
	}
}

func TestUserRepository_IsGranted(t *testing.T) {
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	for ur := range loadUserFixtures(t, suite.repository.user, userPermissionsTestFixtures) {
		for pr := range loadPermissionFixtures(t, suite.repository.permission, ur.given.Permission) {
			add := []*userPermissionsEntity{{
				UserID:              ur.got.ID,
				PermissionSubsystem: pr.got.Subsystem,
				PermissionModule:    pr.got.Module,
				PermissionAction:    pr.got.Action,
			}}
			for _ = range loadUserPermissionsFixtures(t, suite.repository.userPermissions, add) {
				exists, err := suite.repository.user.IsGranted(ur.given.ID, pr.given.Permission())

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

func TestUserRepository_SetPermissions(t *testing.T) {
	t.Skip("not implemented")
}

type userFixtures struct {
	got, given userEntity
}

func loadUserFixtures(t *testing.T, r UserRepository, f []*userEntity) chan userFixtures {
	data := make(chan userFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.Insert(given)
			if err != nil {
				t.Errorf("user cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Logf("user has been created, got id %d", entity.ID)
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
