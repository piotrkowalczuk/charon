package charond

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/ntypes"
)

func untouched(given, created, removed int64) int64 {
	switch {
	case given < 0:
		return -1
	case given == 0:
		return -2
	case given < created:
		return 0
	default:
		return given - created
	}
}

func none() *empty.Empty {
	return &empty.Empty{}
}

func allocNilString(s *ntypes.String) ntypes.String {
	if s == nil {
		return ntypes.String{}
	}
	return *s
}

func allocNilBool(b *ntypes.Bool) ntypes.Bool {
	if b == nil {
		return ntypes.Bool{}
	}
	return *b
}
