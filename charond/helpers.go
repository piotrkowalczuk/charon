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
