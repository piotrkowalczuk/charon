// +build unit,!postgres,!e2e

package main

import "testing"

func TestUntouched(t *testing.T) {
	data := []struct {
		given, created, removed, untouched int64
	}{
		{
			given:     -1000,
			created:   2,
			removed:   10,
			untouched: -1,
		},
		{
			given:     0,
			created:   0,
			removed:   0,
			untouched: -2,
		},
		{
			given:     5,
			created:   3,
			removed:   0,
			untouched: 2,
		},
		{
			given:     10,
			created:   0,
			removed:   0,
			untouched: 10,
		},
		{
			given:     100,
			created:   100,
			removed:   100,
			untouched: 0,
		},
		{
			given:     5,
			created:   0,
			removed:   100,
			untouched: 5,
		},
		{
			given:     5,
			created:   6,
			removed:   0,
			untouched: 0,
		},
	}

	for _, d := range data {
		u := untouched(d.given, d.created, d.removed)
		if u != d.untouched {
			t.Errorf("wrong value, expected %d but got %d", d.untouched, u)
		}
	}
}
