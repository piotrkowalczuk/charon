#!/usr/bin/env bash

SERVICE=charon

cd ./internal/model && charong && cd -
goimports -w ./internal/model
mockery -dir=./internal/model -case=underscore -all -inpkg
mockery -case=underscore -dir=./${SERVICE}rpc -all -output=./${SERVICE}test -outpkg=${SERVICE}test
goimports -w ./internal/model

