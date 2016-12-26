all: build

ARTIFACT = chart
OS = $(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH ?= amd64
GOOS ?= ${OS}
TAG ?= latest

build build-linux build-darwin:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -o ${ARTIFACT} -a .

test:
	go test

run: build
	./${ARTIFACT}

release-linux: GOOS = linux
release-linux: build-linux tar-gzip-linux

release-darwin: GOOS = darwin
release-darwin: build-darwin tar-gzip-darwin

tar-gzip-linux tar-gzip-darwin: BINARY = ${ARTIFACT}-${TAG}-${GOOS}
tar-gzip-linux tar-gzip-darwin:
	tar -cf ${BINARY}.tar ${ARTIFACT}
	gzip -f ${BINARY}.tar
	rm -rf ${ARTIFACT}

release: test release-linux release-darwin
	git tag ${TAG}
