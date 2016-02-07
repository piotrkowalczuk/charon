package main

import (
	"reflect"
	"testing"
)

var (
	permissionTestFixtures = []*permissionEntity{
		{
			Subsystem: "subsystem",
			Module:    "module",
			Action:    "action",
		},
		{
			Module: "module",
			Action: "action",
		},
		{
			Action: "action",
		},
	}
)

func TestPermissionRepository_FindOneByID(t *testing.T) {
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	for res := range loadPermissionFixtures(t, suite.permission, permissionTestFixtures) {
		found, err := suite.permission.FindOneByID(res.got.ID)

		if err != nil {
			t.Errorf("permission cannot be found, unexpected error: %s", err.Error())
			continue
		}

		if assert(t, found != nil, "permission was not found, nil object returned") {
			assertf(t, reflect.DeepEqual(res.got, *found), "created and retrieved entity should be equal, but its not\ncreated: %#v\nfounded: %#v", res.got, found)
		}
	}
}

type permissionFixtures struct {
	got, given permissionEntity
}

func loadPermissionFixtures(t *testing.T, r PermissionRepository, f []*permissionEntity) chan permissionFixtures {
	data := make(chan permissionFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.Insert(given)
			if err != nil {
				t.Errorf("permission cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Logf("permission has been created, got id %d", entity.ID)
			}

			data <- permissionFixtures{
				got:   *entity,
				given: *given,
			}
		}

		close(data)
	}()

	return data
}
