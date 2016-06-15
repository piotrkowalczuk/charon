VERSION=$(shell git describe --tags --always --dirty)
PROTOC=/usr/local/bin/protoc
SERVICE=charon
PACKAGE=github.com/piotrkowalczuk/charon
PACKAGE_TEST=$(PACKAGE)/$(SERVICE)test
PACKAGE_DAEMON=$(PACKAGE)/$(SERVICE)d
PACKAGE_EXAMPLE=$(PACKAGE)/example

PACKAGE_CMD_DAEMON=$(PACKAGE)/cmd/$(SERVICE)d
PACKAGE_CMD_CONTROL=$(PACKAGE)/cmd/$(SERVICE)ctl
PACKAGE_CMD_GENERATOR=$(PACKAGE)/cmd/$(SERVICE)g

BINARY_CMD_DAEMON=.tmp/${SERVICE}d
BINARY_CMD_CONTROL=.tmp/${SERVICE}ctl
BINARY_CMD_GENERATOR=.tmp/${SERVICE}g

#packaging
DIST_PACKAGE_BUILD_DIR=temp
DIST_PACKAGE_DIR=dist
DIST_PACKAGE_TYPE=deb
DIST_PREFIX=/usr
DIST_BINDIR=${DESTDIR}${DIST_PREFIX}/bin

FLAGS=-host=$(CHARON_HOST) \
      	    -port=$(CHARON_PORT) \
      	    -subsystem=$(CHARON_SUBSYSTEM) \
      	    -namespace=$(CHARON_NAMESPACE) \
      	    -test=$(CHARON_TEST) \
      	    -l.format=$(CHARON_LOGGER_FORMAT) \
      	    -l.adapter=$(CHARON_LOGGER_ADAPTER) \
      	    -l.level=$(CHARON_LOGGER_LEVEL) \
      	    -m.engine=$(CHARON_MONITORING_ENGINE) \
      	    -p.address=$(CHARON_POSTGRES_ADDRESS) \
      	    -pwd.strategy=$(CHARON_PASSWORD_STRATEGY) \
      	    -pwd.bcryptcost=$(CHARON_PASSWORD_BCRYPT_COST) \
      	    -mnemo.address=$(CHARON_MNEMOSYNE_ADDRESS)

CMD_TEST=go test -v -coverprofile=profile.out -covermode=atomic

PROTO_PATH=--proto_path=. \
          	    --proto_path=${GOPATH}/src \
          	    --proto_path=${GOPATH}/src/github.com/piotrkowalczuk/mnemosyne \
          	    --proto_path=${GOPATH}/src/github.com/piotrkowalczuk/qtypes \
          	    --proto_path=${GOPATH}/src/github.com/piotrkowalczuk/ntypes \

.PHONY:	all proto rebuild build build-daemon build-control build-example install-generator run test test-short get build package

all: get install

rebuild: install-generator proto gen build

build: build-daemon build-control build-example

build-daemon:
	@go build -o ${BINARY_CMD_DAEMON} ${PACKAGE_CMD_DAEMON}

build-control:
	@go build -o ${BINARY_CMD_CONTROL} ${PACKAGE_CMD_CONTROL}

build-example:
	@go build -o example/client/client ${PACKAGE_EXAMPLE}/client

install-generator:
	@go install ${PACKAGE_CMD_GENERATOR}

gen:
	@go generate ./${SERVICE}d
	@ls -al ${SERVICE}d | grep pqt
	@go generate .
	@ls -al | grep pb.go

run:
	@${BINARY_CMD_DAEMON} ${FLAGS}

test:
	@${CMD_TEST} ${PACKAGE}
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} ${PACKAGE_DAEMON} -p.address=$(CHARON_POSTGRES_ADDRESS)
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} ${PACKAGE_TEST}

test-short:
	@${CMD_TEST} -short ${PACKAGE}
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} -short ${PACKAGE_TEST}

get:
	@go get github.com/Masterminds/glide
	@go get github.com/smartystreets/goconvey/convey
	@glide install

install:
	@go install ${PACKAGE_CMD_DAEMON}
	@go install ${PACKAGE_CMD_CONTROL}
	@go install ${PACKAGE_CMD_GENERATOR}