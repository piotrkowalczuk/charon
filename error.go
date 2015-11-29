package charon

const (
	// ErrDescUserWithIDExists this error should not happen. Id is auto generated.
	ErrDescUserWithIDExists = "charon: user with such id already exists"
	// ErrDescUserWithUsernameExists can be returned if client is trying to create or modify a user,
	// but user with such a username already exists in database.
	ErrDescUserWithUsernameExists = "charon: user with such username already exists"
)
