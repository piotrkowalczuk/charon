#!/usr/bin/env bash

SERVICE=charon
cd ./internal/model && charong && cd -
goimports -w ./internal/model
mockery -case=underscore -dir=./internal/model -all -output=./internal/model/modelmock -outpkg=modelmock
mockery -case=underscore -dir=./internal/session -all -output=./internal/session/sessionmock -outpkg=sessionmock
mockery -case=underscore -dir=./internal/password -all -output=./internal/password/passwordmock -outpkg=passwordmock
mockery -case=underscore -dir=./internal/service -all -output=./internal/service/servicemock -outpkg=servicemock
mockery -case=underscore -dir=./${SERVICE}rpc -all -output=./${SERVICE}test -outpkg=${SERVICE}test
goimports -w ./internal/model

