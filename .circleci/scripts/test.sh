#!/usr/bin/env bash

: ${TEST_RESULTS:=.}

set -e


gotestsum --junitfile results.xml -- -p=1 -count=1 -race -coverprofile=cover-source.out -covermode=atomic -v ./...
cat cover-source.out | grep -v '.pb.go' > cover-step1.out
cat cover-step1.out | grep -v 'mock' > cover-step2.out
cat cover-step2.out | grep -v '/pb/' > cover-step3.out
cat cover-step3.out | grep -v '/example/' > cover-step4.out
cat cover-step4.out | grep -v '.pqt.go' > cover.out

