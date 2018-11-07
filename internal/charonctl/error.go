package charonctl

type Error struct {
	Msg string
	Err error
}

// Error implements error interface.
func (e *Error) Error() string {
	return e.Msg
}
