# Configuration
.DEFAULT_GOAL := help
BINARY ?= company_http
APP_VERSION ?= $(shell git describe --tags --always)

# Tool Versions
GOLANGCI_VERSION ?= v1.50.1
MOCKERY_VERSION ?= v2.14.1

# Main targets
.PHONY: build
build: tidy ## Build the main binary
	go build -o $(BINARY) ./cmd/http

.PHONY: test
test: tidy unittest integration-test ## Run all tests

.PHONY: unittest
unittest: tidy ## Run unit tests
	go test -short -race -v ./...

.PHONY: integration-test
integration-test: ## Run integration tests
	$(MAKE) docker/build
	$(MAKE) docker/run || true
	go clean -testcache
	go test -race -v -timeout 10m ./... -run Integration || true
	$(MAKE) docker/stop

# Docker
.PHONY: docker/build
docker/build: ## Build Docker image for the main binary
	docker build \
		-t $(BINARY):$(APP_VERSION) \
		-t $(BINARY):latest \
		-f ./Dockerfile \
		.

.PHONY: docker/run
docker/run: ## Run the binary in docker
	docker compose up -d

.PHONY: docker/stop
docker/stop: ## Stop & tear down
	docker compose down --volumes

# Maintenance
.PHONY: tidy
tidy: ## Run `go mod tidy`
	go mod tidy

.PHONY: lint
lint: selflint ## Run golangci-lint
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION) run --timeout=5m

.PHONY: generate/mocks
generate/mocks: ## Generate mocks
	go run github.com/vektra/mockery/v2@$(MOCKERY_VERSION) --all --dir ./pkg/domain --output ./pkg/domain/mocks

.PHONY: selflint
selflint: ## Run Makefile linters
	@rm -f .selflint.fail
# check for the ## doc string
	@t=$$(egrep '^[a-zA-Z0-9_/-]+:.*$$' $(MAKEFILE_LIST) | grep -v '\#\#' | grep -v '=' | egrep -v '\# ?nolint') && echo $$t | while read -d: line; do echo -e "selflint: Makefile target '$$line' is missing documentation"; touch .selflint.fail; done || true
# check for .PHONY
	@t=$$(egrep '^[a-zA-Z0-9_/-]+:.*?' $(MAKEFILE_LIST) | egrep -v '\# ?nolint' | egrep -o '^[^:]+' | sort -u); echo $$t | while read -d' ' line; do egrep "^\.PHONY: $$line" $(MAKEFILE_LIST) >/dev/null || ( echo "selflint: Makefile target '$$line' should be .PHONY" && touch .selflint.fail ); done
# fail if any issues are found
	@if [ -f .selflint.fail ]; then echo "Self-lint failed!" && rm -f .selflint.fail && exit 1; fi

.PHONY: help
help: ## Show this help message
	@egrep '^[a-zA-Z0-9_/-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
