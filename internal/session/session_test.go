package session_test

import (
	"testing"

	"github.com/piotrkowalczuk/charon/internal/session"
)

func TestUserIDFromSessionActorID(t *testing.T) {
	success := map[string]int64{
		"charon:user:1":              1,
		"charon:user:0":              0,
		"charon:user:12312412512512": 12312412512512,
	}

	for given, expected := range success {
		userID, err := session.ActorID(given).UserID()
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		} else if userID != expected {
			t.Errorf("wrong user id retrieved from session subject id expected %d, got %d", expected, userID)
		}
	}

	failures := []string{
		"",
		"charon",
		"charon:",
		"charon:user",
		"charon:user:",
		":user:1",
		"user:1",
		":1",
		"1",
		"1231251251241241241251251",
		"charon:resu:52151235125123",
		"charon:u:52151235125123",
		"charon:user:1234567890x",
	}

	for _, given := range failures {
		_, err := session.ActorID(given).UserID()
		if err == nil {
			t.Errorf("expected error %s", err.Error())
		}
	}
}
