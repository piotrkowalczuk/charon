SERVICE=charon
VERSION=$(shell git describe --tags --always --dirty)
ifeq ($(version),)
	TAG=${VERSION}
else
	TAG=$(version)
endif

PACKAGE=github.com/piotrkowalczuk/charon
PACKAGE_CMD_DAEMON=$(PACKAGE)/cmd/$(SERVICE)d
PACKAGE_CMD_CONTROL=$(PACKAGE)/cmd/$(SERVICE)ctl
PACKAGE_CMD_GENERATOR=$(PACKAGE)/cmd/$(SERVICE)g


LDFLAGS = -X 'main.version=$(VERSION)'

.PHONY:	all version build install gen test cover get publish

all: get install

version:
	echo ${VERSION} > publish/VERSION.txt

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -o bin/${SERVICE}g ${PACKAGE_CMD_GENERATOR}
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -o bin/${SERVICE}d ${PACKAGE_CMD_DAEMON}
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -o bin/${SERVICE}ctl ${PACKAGE_CMD_CONTROL}

install:
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_GENERATOR}
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_DAEMON}
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_CONTROL}

gen:
	#./scripts/generate.sh
	bash ./.circleci/scripts/generate.sh golang

test:
	./.circleci/scripts/test.sh
	go tool cover -func=cover.out | tail -n 1

cover: test
	go tool cover -html=cover.out

get:
	get-tool github.com/golang/protobuf/proto v1.2.0
	go get -u gotest.tools/gotestsum
	go get -u github.com/vektra/mockery/cmd/mockery
	go get -u github.com/golang/dep/cmd/dep
	dep ensure


get-tool:
	rm -rf .tmp/tools/$< && \
		mkdir -p .tmp/tools/$< && \
		cd .tmp/tools/$< && \
		GO111MODULE=on go mod init && \
		GO111MODULE=on go get $<@$(word 2,$^)
		GO111MODULE=on go install $<

publish:
	docker build \
		--build-arg VCS_REF=${VCS_REF} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		-t piotrkowalczuk/${SERVICE}:${TAG} .
	docker push piotrkowalczuk/${SERVICE}:${TAG}

setup-python:
	python3 -m venv venv
	source ./venv/bin/activate.fish
	pip install grpc_tools
	pip install grpcio-tools