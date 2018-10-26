all: build

ARTIFACT = chart
OS = $(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH ?= amd64
GOOS ?= ${OS}
TAG ?= latest

build build-linux build-darwin:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 GO111MODULE=on go build -o ${ARTIFACT} -a .

test: ## Run tests
test:
	GO111MODULE=on go test ./...

run: ## Build for default architecture and run immediately.
run: build
	./${ARTIFACT}

release-linux: ## Build a release for linux.
release-linux: GOOS = linux
release-linux: build-linux tar-gzip-linux

release-darwin: ## Build a release for darwin.
release-darwin: GOOS = darwin
release-darwin: build-darwin tar-gzip-darwin

tar-gzip-linux tar-gzip-darwin: BINARY = ${ARTIFACT}-${TAG}-${GOOS}
tar-gzip-linux tar-gzip-darwin:
	tar -cf ${BINARY}.tar ${ARTIFACT}
	gzip -f ${BINARY}.tar
	rm -rf ${ARTIFACT}

release: ## Builds new release and git tags it.
release: test release-linux release-darwin
	git tag ${TAG}


help: ## prints this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
