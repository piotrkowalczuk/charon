PROTOC=/usr/local/bin/protoc
SERVICE=charon
PACKAGE=github.com/piotrkowalczuk/charon
PACKAGE_DAEMON=$(PACKAGE)/$(SERVICE)d
PACKAGE_CONTROLL=$(PACKAGE)/$(SERVICE)ctl
PACKAGE_TEST=$(PACKAGE)/$(SERVICE)test
BINARY_DAEMON=${SERVICE}d/${SERVICE}d
BINARY_CONTROLL=${SERVICE}d/${SERVICE}ctl

FLAGS=-host=$(CHARON_HOST) \
      	    -port=$(CHARON_PORT) \
      	    -subsystem=$(CHARON_SUBSYSTEM) \
      	    -namespace=$(CHARON_NAMESPACE) \
      	    -l.format=$(CHARON_LOGGER_FORMAT) \
      	    -l.adapter=$(CHARON_LOGGER_ADAPTER) \
      	    -l.level=$(CHARON_LOGGER_LEVEL) \
      	    -m.engine=$(CHARON_MONITORING_ENGINE) \
      	    -ps.connectionstring=$(CHARON_POSTGRES_CONNECTION_STRING) \
      	    -ps.retry=$(CHARON_POSTGRES_RETRY) \
      	    -pwd.strategy=$(CHARON_PASSWORD_STRATEGY) \
      	    -pwd.bcryptcost=$(CHARON_PASSWORD_BCRYPT_COST) \
      	    -mnemo.address=$(CHARON_MNEMOSYNE_ADDRESS)

.PHONY:	all proto build build-daemon run test test-unit test-postgres

all: proto build test run

proto:
	@${PROTOC} --proto_path=. \
	    --proto_path=${GOPATH}/src/github.com/piotrkowalczuk/mnemosyne \
	    --proto_path=${GOPATH}/src/github.com/piotrkowalczuk/protot \
	    --proto_path=${GOPATH}/src/github.com/piotrkowalczuk/nilt \
	    --go_out=Mmnemosyne.proto=github.com/piotrkowalczuk/mnemosyne,Mprotot.proto=github.com/piotrkowalczuk/protot,Mnilt.proto=github.com/piotrkowalczuk/nilt,plugins=grpc:. \
		${SERVICE}.proto
	@ls -al | grep pb.go

build: build-daemon build-controll

build-daemon:
	@go build -o ${BINARY_DAEMON} ${PACKAGE_DAEMON}

build-controll:
	@go build -o ${BINARY_CONTROLL} ${PACKAGE_CONTROLL}

run:
	@${BINARY_DAEMON} ${FLAGS}

test: test-unit test-postgres

test-unit:
	@go test -v ${PACKAGE}
	@go test -v ${PACKAGE_TEST}
	@go test -v -tags=unit ${PACKAGE_DAEMON}

test-postgres:
	@go test -v -tags=postgres ${PACKAGE_DAEMON} ${FLAGS}

get:
	@go get ${PACKAGE}
	@go get ${PACKAGE_TEST}
	@go get ${PACKAGE_DAEMON}