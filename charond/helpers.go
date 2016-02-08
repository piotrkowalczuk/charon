package main

import "github.com/piotrkowalczuk/nilt"

func nilString(ns *nilt.String) nilt.String {
	if ns == nil {
		return nilt.String{}
	}

	return *ns
}

func nilBool(nb *nilt.Bool) nilt.Bool {
	if nb == nil {
		return nilt.Bool{}
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
