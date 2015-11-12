PROTOC=/usr/local/bin/protoc
SERVICE=charon
PACKAGE=github.com/piotrkowalczuk/charon
PACKAGE_DAEMON=$(PACKAGE)/$(SERVICE)d
BINARY=${SERVICE}d/${SERVICE}d

FLAGS=

.PHONY:	all proto build build-daemon run test test-unit test-postgres

all: proto build test run

proto:
	@${PROTOC} --proto_path=${GOPATH}/src \
	    --proto_path=. \
	    --go_out=plugins=grpc:. \
	    ${SERVICE}.proto
	@ls -al | grep pb.go

build: build-daemon

build-daemon:
	@go build -o ${BINARY} ${PACKAGE_DAEMON}

run:
	@${BINARY} ${FLAGS}

test: test-unit test-postgres

test-unit:
	@go test -v ${PACKAGE_DAEMON} ${FLAGS}

test-postgres:
	@go test -tags postgres -v ${PACKAGE_DAEMON} ${FLAGS}

get:
	@go get ${PACKAGE_DAEMON}