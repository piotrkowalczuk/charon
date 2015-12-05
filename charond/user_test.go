// +build postgres !unit

package main

import "testing"

func TestUserRepository_Create(t *testing.T) {
	config.parse()

	postgres := initPostgres(config.postgres.connectionString, 0, &testLogger{T: t})
	repo := newUserRepository(postgres)

	success := []struct {
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

TestLoop:
	for _, data := range success {
		entity, err := repo.Create(
			data.username,
			data.password,
			data.firstName,
			data.lastName,
			data.confirmationToken,
			data.isSuperuser,
			data.isStaff,
			data.isActive,
			data.isConfirmed,
		)

		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			continue TestLoop
		}

		if entity.CreatedAt == nil || entity.CreatedAt.IsZero() {
			t.Errorf("invalid created at field, expected valid time but got %v", entity.CreatedAt)
		} else {
			t.Logf("user has been properly created at %v", entity.CreatedAt)
		}

		if entity.Username != data.username {
			t.Errorf("wrong username, expected %s but got %s", data.username, entity.Username)
		}
	}

	if err := tearDownDatabase(postgres); err != nil {
		t.Errorf("unexpected error during database teardown: %s", err.Error())
	}

	postgres.Close()
}

type testLogger struct {
	*testing.T
}

func (tl *testLogger) Log(args ...interface{}) error {
	tl.T.Log(args...)

	return nil
}
