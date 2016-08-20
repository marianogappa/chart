ARTIFACT = chart
OS = $(shell uname | tr [:upper:] [:lower:])

all: build

build: GOOS ?= ${OS}
build: GOARCH ?= amd64
build: clean
		GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -o ${ARTIFACT} -a *.go

clean:
		rm -f ${ARTIFACT}

test:
		go test

run: build
	./${ARTIFACT}
