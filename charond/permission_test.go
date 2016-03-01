// +build postgres !unit

package main

import (
	"reflect"
	"testing"

	"github.com/piotrkowalczuk/charon"
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

	for res := range loadPermissionFixtures(t, suite.repository.permission, permissionTestFixtures) {
		found, err := suite.repository.permission.FindOneByID(res.got.ID)

		if err != nil {
			t.Errorf("permission cannot be found, unexpected error: %s", err.Error())
			continue
		}

		if assert(t, found != nil, "permission was not found, nil object returned") {
			assertf(t, reflect.DeepEqual(res.got, *found), "created and retrieved entity should be equal, but its not\ncreated: %#v\nfounded: %#v", res.got, found)
		}
	}
}

func TestPermissionRepository_Register(t *testing.T) {
	suite := setupPostgresSuite(t)
	defer suite.teardown(t)

	data := []struct {
		created, removed, untouched int64
		permissions                 charon.Permissions
	}{
		{
			created:     int64(len(charon.AllPermissions)),
			permissions: charon.AllPermissions,
		},
		{
			untouched:   int64(len(charon.AllPermissions)),
			permissions: charon.AllPermissions,
		},
		{
			untouched: 1,
			removed:   int64(len(charon.AllPermissions) - 1),
			permissions: charon.Permissions{
				charon.UserCanCreate,
			},
		},
		{
			created: 1,
			removed: 1,
			permissions: charon.Permissions{
				charon.Permission("charon:fakemodule:fakeaction"),
			},
		},
		{
			created: 1,
			permissions: charon.Permissions{
				charon.Permission("fakesystem:fakemodule:fakeaction"),
			},
		},
		{
			removed:     1,
			created:     int64(len(charon.AllPermissions)),
			permissions: charon.AllPermissions,
		},
	}

	for i, d := range data {
		created, untouched, removed, err := suite.repository.permission.Register(d.permissions)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		if created != d.created {
			t.Errorf("expected different number of created permissions, expected %d got %d for set %d", d.created, created, i)
		}
		if untouched != d.untouched {
			t.Errorf("expected different number of untouched permissions, expected %d got %d for set %d", d.untouched, untouched, i)
		}
		if removed != d.removed {
			t.Errorf("expected different number of removed permissions, expected %d got %d for set %d", d.removed, removed, i)
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
