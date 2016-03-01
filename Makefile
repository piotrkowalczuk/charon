PROTOC=/usr/local/bin/protoc
SERVICE=charon
PACKAGE=github.com/piotrkowalczuk/charon
PACKAGE_DAEMON=$(PACKAGE)/$(SERVICE)d
PACKAGE_CONTROL=$(PACKAGE)/$(SERVICE)ctl
PACKAGE_GENERATOR=$(PACKAGE)/$(SERVICE)g
PACKAGE_EXAMPLE=$(PACKAGE)/example
PACKAGE_TEST=$(PACKAGE)/$(SERVICE)test

BINARY_DAEMON=${SERVICE}d/${SERVICE}d
BINARY_CONTROL=${SERVICE}ctl/${SERVICE}ctl
BINARY_GENERATOR=${SERVICE}g/${SERVICE}g

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
      	    -l.format=$(CHARON_LOGGER_FORMAT) \
      	    -l.adapter=$(CHARON_LOGGER_ADAPTER) \
      	    -l.level=$(CHARON_LOGGER_LEVEL) \
      	    -m.engine=$(CHARON_MONITORING_ENGINE) \
      	    -ps.connectionstring=$(CHARON_POSTGRES_CONNECTION_STRING) \
      	    -pwd.strategy=$(CHARON_PASSWORD_STRATEGY) \
      	    -pwd.bcryptcost=$(CHARON_PASSWORD_BCRYPT_COST) \
      	    -mnemo.address=$(CHARON_MNEMOSYNE_ADDRESS)

CMD_TEST=go test -v -coverprofile=profile.out -covermode=atomic

PROTO_PATH=--proto_path=. \
          	    --proto_path=${GOPATH}/src/github.com/piotrkowalczuk/mnemosyne \
          	    --proto_path=${GOPATH}/src/github.com/piotrkowalczuk/protot \
          	    --proto_path=${GOPATH}/src/github.com/piotrkowalczuk/nilt \

.PHONY:	all proto rebuild build build-daemon build-control build-example install-generator run test test-unit test-postgres get buld package

all: rebuild test run

proto:
	@${PROTOC} ${PROTO_PATH}--go_out=Mmnemosyne.proto=github.com/piotrkowalczuk/mnemosyne,Mprotot.proto=github.com/piotrkowalczuk/protot,Mnilt.proto=github.com/piotrkowalczuk/nilt,plugins=grpc:. \
		${SERVICE}.proto
	@ls -al | grep pb.go

rebuild: install-generator proto generate build

build: build-daemon build-control build-example

build-daemon:
	@go build -o ${BINARY_DAEMON} ${PACKAGE_DAEMON}

build-control:
	@go build -o ${BINARY_CONTROL} ${PACKAGE_CONTROL}

build-example:
	@go build -o example/client/client ${PACKAGE_EXAMPLE}/client

install-generator:
	@go install ${PACKAGE_GENERATOR}

generate:
	@go generate ./...
	@goimports -w ./charond/schema.go

run:
	@${BINARY_DAEMON} ${FLAGS}

test: test-unit test-postgres

test-unit:
	@${CMD_TEST} ${PACKAGE}
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} ${PACKAGE_TEST}
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} -tags=unit ${PACKAGE_DAEMON}
	@cat profile.out >> coverage.txt && rm profile.out

test-postgres:
	@${CMD_TEST} -tags=postgres ${PACKAGE_DAEMON} ${FLAGS}
	@cat profile.out >> coverage.txt && rm profile.out

get:
	@go get github.com/smartystreets/goconvey/convey
	@go get ${PACKAGE}
	@go get ${PACKAGE_TEST}
	@go get ${PACKAGE_DAEMON}

install: build
	#install binary
	install -Dm 755 ${BINARY_DAEMON} ${DIST_BINDIR}/${SERVICE}d
	install -Dm 755 ${BINARY_CONTROL} ${DIST_BINDIR}/${SERVICE}ctl
	#install config file
	install -Dm 644 scripts/${SERVICE}.env ${DESTDIR}/etc/${SERVICE}.env
	install -Dm 644 scripts/${SERVICE}.env ${DESTDIR}/etc/${SERVICE}.env
	#install init script
	install -Dm 644 scripts/${SERVICE}.service ${DESTDIR}/etc/systemd/system/${SERVICE}.service

package:
	# export DIST_PACKAGE_TYPE to vary package type (e.g. deb, tar, rpm)
	@if [ -z "$(shell which fpm 2>/dev/null)" ]; then \
		echo "error:\nPackaging requires effing package manager (fpm) to run.\nsee https://github.com/jordansissel/fpm\n"; \
		exit 1; \
	fi

	#run make install against the packaging dir
	mkdir -p ${DIST_PACKAGE_BUILD_DIR} && $(MAKE) install DESTDIR=${DIST_PACKAGE_BUILD_DIR}

	#clean
	mkdir -p ${DIST_PACKAGE_DIR} && rm -f ${DIST_PACKAGE_DIR}/*.${DIST_PACKAGE_TYPE}

	#build package
	fpm --rpm-os linux \
		-s dir \
		-p dist \
		-t ${DIST_PACKAGE_TYPE} \
		-n ${SERVICE} \
		-v `${DIST_PACKAGE_BUILD_DIR}${DIST_PREFIX}/bin/${SERVICE}d -version` \
		--config-files /etc/${SERVICE}.env \
		--config-files /etc/systemd/system/${SERVICE}.service \
		-C ${DIST_PACKAGE_BUILD_DIR} .
