package charond

import (
	"context"
	"database/sql"

	"github.com/piotrkowalczuk/charon/internal/model"
	"github.com/piotrkowalczuk/charon/internal/password"
)

type repositories struct {
	user             model.UserProvider
	userGroups       model.UserGroupsProvider
	userPermissions  model.UserPermissionsProvider
	permission       model.PermissionProvider
	group            model.GroupProvider
	groupPermissions model.GroupPermissionsProvider
	refreshToken     model.RefreshTokenProvider
}

func newRepositories(db *sql.DB) repositories {
	return repositories{
		user:             model.NewUserRepository(db),
		userGroups:       model.NewUserGroupsRepository(db),
		userPermissions:  model.NewUserPermissionsRepository(db),
		permission:       model.NewPermissionRepository(db),
		group:            model.NewGroupRepository(db),
		groupPermissions: model.NewGroupPermissionsRepository(db),
		refreshToken:     model.NewRefreshTokenRepository(db),
	}
}

func execQueries(db *sql.DB, queries ...string) (err error) {
	exec := func(query string) {
		if err != nil {
			return
		}

		_, err = db.Exec(query)
	}

	for _, q := range queries {
		exec(q)
	}

	return
}

func setupDatabase(db *sql.DB) error {
	return execQueries(
		db,
		model.SQL,
	)
}

func teardownDatabase(db *sql.DB) error {
	return execQueries(
		db,
		`DROP SCHEMA IF EXISTS charon CASCADE`,
	)
}

func createDummyTestUser(ctx context.Context, usrProvider model.UserProvider, rftProvider model.RefreshTokenProvider, hasher password.Hasher) error {
	pass, err := hasher.Hash([]byte("test"))
	if err != nil {
		return err
	}
	usr, err := usrProvider.CreateSuperuser(ctx, "test", pass, "Test", "Test")
	if err != nil {
		return err
	}

	_, err = rftProvider.Create(ctx, &model.RefreshTokenEntity{Token: "test", UserID: usr.ID})
	return err
}
