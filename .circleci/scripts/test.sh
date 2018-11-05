#!/usr/bin/env bash

: ${TEST_RESULTS:=.}

set -e

gotestsum --junitfile results.xml -- -count=1 -race -coverprofile=cover-source.out -covermode=atomic -v ./...
cat cover-source.out | grep -v '.pb.go' > cover-step1.out
cat cover-step1.out | grep -v '.pqt.go' > cover.out

