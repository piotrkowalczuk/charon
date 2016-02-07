// +build postgres !unit

package main

import "testing"

var (
	groupTestFixtures = []*groupEntity{}
)

type groupFixtures struct {
	got, given groupEntity
}

func loadGroupFixtures(t *testing.T, r GroupRepository, f []*groupEntity) chan groupFixtures {
	data := make(chan groupFixtures, 1)

	go func() {
		for _, given := range f {
			entity, err := r.Insert(given)
			if err != nil {
				t.Errorf("group cannot be created, unexpected error: %s", err.Error())
				continue
			} else {
				t.Logf("group has been created, got id %d", entity.ID)
			}

			data <- groupFixtures{
				got:   *entity,
				given: *given,
			}
		}

		close(data)
	}()

	return data
}
