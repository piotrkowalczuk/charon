#!/usr/bin/env bash

SERVICE=charon
PROTO_INCLUDE="-I=./${SERVICE}rpc -I=/usr/include -I=./vendor/github.com/piotrkowalczuk"
PROTO_FILES="${SERVICE}rpc/*.proto"


cd ./internal/model && charong && cd -
goimports -w ./internal/model
mockery -dir=./internal/model -case=underscore -all -inpkg
goimports -w ./internal/model


protoc ${PROTO_INCLUDE} --go_out=plugins=grpc:${GOPATH}/src ${PROTO_FILES}
python -m grpc_tools.protoc ${PROTO_INCLUDE} --python_out=./${SERVICE}rpc --grpc_python_out=./${SERVICE}rpc ${PROTO_FILES}
goimports -w ./${SERVICE}rpc

mockery -case=underscore -dir=./${SERVICE}rpc -all -output=./${SERVICE}test -outpkg=${SERVICE}test

ls -lha ./${SERVICE}rpc | grep pb.go
ls -lha ./${SERVICE}rpc | grep pb2.py
ls -lha ./${SERVICE}rpc | grep grpc.py

