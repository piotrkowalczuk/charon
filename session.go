package charon

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	subjectIDPrefix = "charon:user:"
)

// SubjectID is globally unique identifier that in format "charon:user:<user_id>".
type SubjectID string

// SubjectIDFromInt64 allocate SessionSubjectID using given user id.
func SubjectIDFromInt64(userID int64) SubjectID {
	return SubjectID(subjectIDPrefix + strconv.FormatInt(userID, 10))
}

// String implements fmt.Stringer interface.
func (ssi SubjectID) String() string {
	return string(ssi)
}

// UserID returns user id if possible, otherwise an error.
func (ssi SubjectID) UserID() (int64, error) {
	if len(ssi) < 13 {
		return 0, errors.New("charon: session subject id to short, min length 13 characters")
	}
	if ssi[:12] != subjectIDPrefix {
		return 0, fmt.Errorf("charon: session subject id wrong prefix expected %s, got %s", subjectIDPrefix, ssi[:12])
	}

	return strconv.ParseInt(string(ssi)[12:], 10, 64)
}
