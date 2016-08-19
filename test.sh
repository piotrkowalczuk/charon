#!/usr/bin/env bash

export CHARON_POSTGRES_ADDRESS="postgres://localhost/test?sslmode=disable"

make "$@"
