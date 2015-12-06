// +build postgres !unit

package main

import (
	"testing"

	"github.com/piotrkowalczuk/nilt"
)

func TestUserRepository_Create(t *testing.T) {
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for res := range generateUserRepositoryData(t, suite) {
		if res.got.CreatedAt == nil || res.got.CreatedAt.IsZero() {
			t.Errorf("invalid created at field, expected valid time but got %v", res.got.CreatedAt)
		} else {
			t.Logf("user has been properly created at %v", res.got.CreatedAt)
		}

		if res.got.Username != res.given.username {
			t.Errorf("wrong username, expected %s but got %s", res.given.username, res.got.Username)
		}
	}
}

func TestUserRepository_UpdateOneByID(t *testing.T) {
	suffix := "_modified"
	suite := &postgresSuite{}
	suite.setup(t)
	defer suite.teardown(t)

	for res := range generateUserRepositoryData(t, suite) {
		user := res.got
		modified, err := suite.user.UpdateOneByID(
			user.ID,
			&nilt.String{String: user.Username + suffix, Valid: true},
			nil,
			&nilt.String{String: user.FirstName + suffix, Valid: true},
			&nilt.String{String: user.LastName + suffix, Valid: true},
			&nilt.Bool{Bool: true, Valid: true},
			&nilt.Bool{Bool: true, Valid: true},
			&nilt.Bool{Bool: true, Valid: true},
			&nilt.Bool{Bool: true, Valid: true},
		)

		if err != nil {
			t.Errorf("user cannot be modified, unexpected error: %s", err.Error())
			continue
		} else {
			t.Logf("user has been modified with id: %d", modified.ID)
		}

		if assertf(t, modified.UpdatedAt != nil && !modified.UpdatedAt.IsZero(), "invalid updated at field, expected valid time but got %v", modified.UpdatedAt) {
			t.Logf("user has been properly modified at %v", modified.UpdatedAt)
		}
		assertf(t, modified.Username == user.Username+suffix, "wrong username, expected %s but got %s", user.Username+suffix, modified.Username)
		assertf(t, modified.Password == user.Password, "wrong password, expected %s but got %s", user.Password, modified.Password)
		assertf(t, modified.FirstName == user.FirstName+suffix, "wrong first name, expected %s but got %s", user.FirstName+suffix, modified.FirstName)
		assertf(t, modified.LastName == user.LastName+suffix, "wrong last name, expected %s but got %s", user.LastName+suffix, modified.LastName)
		assert(t, modified.IsSuperuser, "user should become a superuser")
		assert(t, modified.IsActive, "user should be active")
		assert(t, modified.IsConfirmed, "user should be confirmed")
		assert(t, modified.IsStaff, "user should become a staff")
	}
}

type userRepositoryFixture struct {
	got   userEntity
	given struct {
		username, password, firstName, lastName, confirmationToken string
		isSuperuser, isStaff, isActive, isConfirmed                bool
	}
}

func generateUserRepositoryData(t *testing.T, suite *postgresSuite) chan userRepositoryFixture {
	given := []struct {
		username, password, firstName, lastName, confirmationToken string
		isSuperuser, isStaff, isActive, isConfirmed                bool
	}{
		{
			username:          "johnsnow@gmail.com",
			password:          "secret",
			firstName:         "John",
			lastName:          "Snow",
			confirmationToken: "1234567890",
		},
		{
			username:          "1",
			password:          "2",
			firstName:         "3",
			lastName:          "4",
			confirmationToken: "5",
		},
	}
	data := make(chan userRepositoryFixture, 1)

	go func() {
		for _, g := range given {
			entity, err := suite.user.Create(
				g.username,
				g.password,
				g.firstName,
				g.lastName,
				g.confirmationToken,
				g.isSuperuser,
				g.isStaff,
				g.isActive,
				g.isConfirmed,
			)
			if err != nil {
				t.Errorf("user cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Logf("user has been created with id: %d", entity.ID)
			}

			data <- userRepositoryFixture{
				got:   *entity,
				given: g,
			}
		}

		close(data)
	}()

	return data
}
