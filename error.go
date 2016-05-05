package charon

// TODO: remove?

const (
	// ErrDescUserWithIDExists this error should not happen. Id is auto generated.
	ErrDescUserWithIDExists = "charon: user with such id already exists"
	// ErrDescUserWithUsernameExists can be returned if client is trying to create or modify a user,
	// but user with such a username already exists in database.
	ErrDescUserWithUsernameExists = "charon: user with such username already exists"
)

const (
	// ErrDescGroupWithIDExists this error should not happen. Id is auto generated.
	ErrDescGroupWithIDExists = "charon: group with such id already exists"
	// ErrDescGroupWithNameExists can be returned if client is trying to create or modify a group,
	// but group with such a name already exists in database.
	ErrDescGroupWithNameExists = "charon: group with such name already exists"
)
