package charon

import (
	"errors"
	"strconv"

	"github.com/piotrkowalczuk/mnemosyne"
)

const (
	sessionKeyUserID = "charon_user_id"
)

// UserIDFromSession returns int64 value if it exists under expected key.
func UserIDFromSession(s *mnemosyne.Session) (int64, error) {
	if s.Bag == nil {
		return 0, errors.New("does not exists")
	}

	return strconv.ParseInt(s.Bag[sessionKeyUserID], 10, 64)
}
