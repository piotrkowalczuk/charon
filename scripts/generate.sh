#!/usr/bin/env bash

SERVICE=charon
PROTO_INCLUDE="-I=./${SERVICE}rpc -I=/usr/include -I=./vendor/github.com/piotrkowalczuk"
PROTO_FILES="${SERVICE}rpc/*.proto"


cd ./internal/model && charong && cd -
mockery -dir=./internal/model -all -inpkg -output_file=./internal/model/mocks.go
goimports -w ./internal/model/schema.pqt.go ./internal/model/mocks.go


protoc ${PROTO_INCLUDE} --go_out=plugins=grpc:${GOPATH}/src ${PROTO_FILES}
python -m grpc_tools.protoc ${PROTO_INCLUDE} --python_out=./${SERVICE}rpc --grpc_python_out=./${SERVICE}rpc ${PROTO_FILES}
goimports -w ./${SERVICE}rpc

mockery -case=underscore -dir=./${SERVICE}rpc -name=.*Client -output=./${SERVICE}test -output_file=${SERVICE}rpc.mock.go -output_pkg_name=${SERVICE}test

ls -lha ./${SERVICE}rpc | grep pb.go
ls -lha ./${SERVICE}rpc | grep pb2.py
ls -lha ./${SERVICE}rpc | grep grpc.py

