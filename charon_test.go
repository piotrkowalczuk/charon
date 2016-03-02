package charon

import "testing"

func TestUser_Name(t *testing.T) {
	given := &User{FirstName: "John", LastName: "Snow"}
	expected := "John Snow"

	got := given.Name()
	if got != expected {
		t.Errorf("output do not match, expected %s but got %s", expected, got)
	}
}
