#!/usr/bin/env bash

rm coverage.txt
set -e
echo "mode: atomic" > coverage.txt

for d in $(go list ./... | grep -v /vendor | grep -v /example | grep -v /charonrpc | grep -v /charontest); do
    if [ "$d" == "github.com/piotrkowalczuk/charon/charond" ]; then
        COVER_PKG="github.com/piotrkowalczuk/charon/charond,github.com/piotrkowalczuk/charon/internal/model"
    else
        COVER_PKG=$d
    fi

	go test -race -coverprofile=profile.out -coverpkg=$COVER_PKG -covermode=atomic $d
	if [ -f profile.out ]; then
		tail -n +2 profile.out >> coverage.txt
		rm profile.out
	fi
done