GOPATH ?= $(HOME)/go
BIN_DIR = $(GOPATH)/bin
TMPDIR ?= $(shell dirname $$(mktemp -u))

# Project specific variables

PACKAGE = candles
NAMESPACE = github.com/$(PACKAGE)
COVER_FILE ?= $(TMPDIR)/$(PACKAGE)-coverage.out

# Tools

GOLANG_CI_LINT = golangci-lint

$(GOLANG_CI_LINT):
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.25.1

.PHONY: tools
tools: $(GOLANG_CI_LINT) ## Install all needed tools, e.g. for static checks

# Main targets
all: test build
.DEFAULT_GOAL := all

.PHONY: build
build: ## Build the project binary
	go build ./cmd/$(PACKAGE)/

.PHONY: build_race
build_race: ## Build the project binary with data race detector
	go build -race ./cmd/$(PACKAGE)/

.PHONY: test
test: ## Run unit (short) tests
	go test -race ./... -coverprofile=$(COVER_FILE)
	go tool cover -func=$(COVER_FILE) | grep ^total

$(COVER_FILE):
	$(MAKE) test

.PHONY: cover
cover: $(COVER_FILE) ## Output coverage in human readable form in html
	go tool cover -html=$(COVER_FILE)
	rm -f $(COVER_FILE)

.PHONY: lint
lint: $(GOLANG_CI_LINT) ## Check the project with lint
	golangci-lint run

.PHONY: static_check
static_check: lint ## Run static checks all over the project

.PHONY: check
check: static_check test ## Check project with static checks and unit tests

.PHONY: run
run: build ## Start the project
	./$(PACKAGE)

.PHONY: run_race
run_race: build_race ## Start the project with data race detector
	./$(PACKAGE)

.PHONY: clean
clean: ## Clean the project from built files
	rm -f ./$(PACKAGE) $(COVER_FILE)

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: dependencies-download
dependencies-download:
	go mod download
