package charon

import (
	"strconv"

	"github.com/piotrkowalczuk/mnemosyne"
)

const (
	sessionKeyUserID = "charon_user_id"
)

// UserIDFromSession returns int64 value if it exists under expected key.
func UserIDFromSession(s *mnemosyne.Session) (int64, error) {
	return strconv.ParseInt(s.Value(sessionKeyUserID), 10, 64)
}
