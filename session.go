package charon

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	sessionSubjectIDPrefix = "charon:user:"
)

// SessionSubjectID is globally unique identifier that in format "charon:user:<user_id>".
type SessionSubjectID string

// NewSessionSubjectID allocate SessionSubjectID using given user id.
func NewSessionSubjectID(userID int64) SessionSubjectID {
	return SessionSubjectID(sessionSubjectIDPrefix + strconv.FormatInt(userID, 10))
}

// String implements fmt.Stringer interface.
func (ssi SessionSubjectID) String() string {
	return string(ssi)
}

// UserID returns user id if possible, otherwise an error.
func (ssi SessionSubjectID) UserID() (int64, error) {
	if len(ssi) < 13 {
		return 0, errors.New("charon: session subject id to short, min length 13 characters")
	}
	if ssi[:12] != sessionSubjectIDPrefix {
		return 0, fmt.Errorf("charon: session subject id wrong prefix expected %s, got %s", sessionSubjectIDPrefix, ssi[:12])
	}

	return strconv.ParseInt(string(ssi)[12:], 10, 64)
}
