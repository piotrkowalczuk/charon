package charon

import (
	"strconv"

	"github.com/piotrkowalczuk/ntypes"
)

func nilString(ns *ntypes.String) ntypes.String {
	if ns == nil {
		return ntypes.String{}
	}

	return *ns
}

func nilBool(nb *ntypes.Bool) ntypes.Bool {
	if nb == nil {
		return ntypes.Bool{}
	}

	return *nb
}

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

func address(host string, port int) string {
	return host + ":" + strconv.FormatInt(int64(port), 10)
}
