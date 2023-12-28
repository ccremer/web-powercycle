# Set Shell to bash, otherwise some targets fail with dash/zsh etc.
SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c

.DEFAULT_GOAL := help

##
## These are some common variables for Make
##

PROJECT_ROOT_DIR = .
PROJECT_NAME ?= web-powercycle
PROJECT_OWNER ?= ccremer

WORK_DIR = $(PWD)/.work

## BUILD:go
BIN_FILENAME ?= $(PROJECT_NAME)
go_bin ?= $(WORK_DIR)/bin
$(go_bin):
	@mkdir -p $@


# extensible array of targets. Modules can add target to this variable for the all-in-one target.
clean_targets := build-clean release-clean
test_targets := test-unit

.PHONY: help
help: ## Show this help
	@grep -E -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: build-bin

.PHONY: build-bin
build-bin: export CGO_ENABLED = 0
build-bin: fmt vet ## Build binary
	@go build $(go_build_args) -o $(BIN_FILENAME) .

build-clean: ## Deletes build artifacts
	rm -rf $(BIN_FILENAME)

.PHONY: test
test: $(test_targets) ## All-in-one test

.PHONY: test-unit
test-unit: ## Run unit tests against code
	go test -race -covermode atomic ./...

.PHONY: fmt
fmt: ## Run 'go fmt' against code
	go fmt ./...

.PHONY: vet
vet: ## Run 'go vet' against code
	go vet ./...

.PHONY: lint
lint: lint-go git-diff ## All-in-one linting

.PHONY: lint-go
lint-go: fmt vet generate ## Run linting for Go code

.PHONY: git-diff
git-diff:
	@echo 'Check for uncommitted changes ...'
	git diff --exit-code

.PHONY: generate
generate: generate-go ## All-in-one code generation

.PHONY: generate-go
generate-go: ## Generate Go artifacts
	@go generate ./...

.PHONY: release-prepare
release-prepare: ## Prepares artifacts for releases

.PHONY: release-clean
release-clean:

.PHONY: clean
clean: $(clean_targets) ## All-in-one target to cleanup local artifacts
