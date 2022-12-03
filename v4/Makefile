
ARTIFACT = chart

build: ## build for default architecture.
build:
	go build -o ${ARTIFACT} -a .

test: ## Run tests
test:
	go test ./...

run: ## Build for default architecture and run immediately.
run: build
	./${ARTIFACT}

release: ## Build, package, tag and publish releases with goreleaser.
release:
	goreleaser --rm-dist

release-dry: ## Build and package releases with goreleaser without publishing or tagging
release-dry:
	goreleaser --snapshot --skip-publish --rm-dist

help: ## prints this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
