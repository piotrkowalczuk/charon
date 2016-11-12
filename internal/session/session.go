package session

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	actorIDPrefix = "charon:user:"
)

// ActorID is globally unique identifier that in format "charon:user:<user_id>".
type ActorID string

// ActorIDFromInt64 allocate ActorID using given user id.
func ActorIDFromInt64(userID int64) ActorID {
	return ActorID(actorIDPrefix + strconv.FormatInt(userID, 10))
}

// String implements fmt.Stringer interface.
func (ai ActorID) String() string {
	return string(ai)
}

// UserID returns user id if possible, otherwise an error.
func (ai ActorID) UserID() (int64, error) {
	if len(ai) < 13 {
		return 0, errors.New("charon: session actor id to short, min length 13 characters")
	}
	if ai[:12] != actorIDPrefix {
		return 0, fmt.Errorf("charon: session actor id wrong prefix expected %s, got %s", actorIDPrefix, ai[:12])
	}

	return strconv.ParseInt(string(ai)[12:], 10, 64)
}
