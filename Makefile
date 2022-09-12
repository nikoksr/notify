export GO111MODULE := on
export GOPROXY = https://proxy.golang.org,direct

###############################################################################
# DEPENDENCIES
###############################################################################

# Install all the build and lint dependencies
setup:
	go mod tidy
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/daixiang0/gci@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/vektra/mockery/v2@latest
.PHONY: setup

###############################################################################
# TESTS
###############################################################################

# Run all the tests
test:
	go test -failfast -race -timeout=5m ./...
.PHONY: test

cover:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
.PHONY: cover

mock:
	go generate ./...
.PHONY: mock

###############################################################################
# CODE HEALTH
###############################################################################

fmt:
	@gofumpt -w -l .

	@gci write --section Standard --section Default --section "Prefix(github.com/nikoksr/notify)" .
.PHONY: fmt


lint:
	@golangci-lint run --config .golangci.yml
.PHONY: lint

ci: lint test
.PHONY: ci

###############################################################################

.DEFAULT_GOAL := ci
