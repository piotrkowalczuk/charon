package main

import "testing"

var (
	permissionTestFixtures = []*permissionEntity{}
)

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
