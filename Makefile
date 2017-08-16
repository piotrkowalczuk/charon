VERSION?=$(shell git describe --tags --always --dirty)
SERVICE=charon

PACKAGE=github.com/piotrkowalczuk/charon
PACKAGE_CMD_DAEMON=$(PACKAGE)/cmd/$(SERVICE)d
PACKAGE_CMD_CONTROL=$(PACKAGE)/cmd/$(SERVICE)ctl
PACKAGE_CMD_GENERATOR=$(PACKAGE)/cmd/$(SERVICE)g


LDFLAGS = -X 'main.version=$(VERSION)'

.PHONY:	all version build install gen test cover get publish

all: get install

version:
	echo ${VERSION} > VERSION.txt

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -o bin/${SERVICE}g ${PACKAGE_CMD_GENERATOR}
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -o bin/${SERVICE}d ${PACKAGE_CMD_DAEMON}
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -o bin/${SERVICE}ctl ${PACKAGE_CMD_CONTROL}

install:
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_GENERATOR}
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_DAEMON}
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_CONTROL}

gen:
	./scripts/generate.sh

test:
	./scripts/test.sh
	go tool cover -func=coverage.txt | tail -n 1

cover: test
	go tool cover -html=coverage.txt

get:
	go get github.com/vektra/mockery/cmd/mockery
	go get github.com/Masterminds/glide
	glide install -v
	cd ./cmd/charong && glide install -v

publish:
	docker build \
		--build-arg VCS_REF=${VCS_REF} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		-t piotrkowalczuk/${SERVICE}:${VERSION} .
	docker push piotrkowalczuk/${SERVICE}:${VERSION}