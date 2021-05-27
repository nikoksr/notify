export GO111MODULE := on
export GOPROXY = https://proxy.golang.org,direct

###############################################################################
# DEPENDENCIES
###############################################################################

# Install all the build and lint dependencies
setup:
	go mod tidy
.PHONY: setup

###############################################################################
# TESTS
###############################################################################

# Run all the tests
test:
	go test -failfast -race -timeout=5m ./...
.PHONY: test

###############################################################################
# CODE HEALTH
###############################################################################

# gofumpt and gci all go files
fmt:
	gofumpt -w .
	gci -w -local github.com/nikoksr/notify .
.PHONY: fmt

# Run all the linters
lint:
	golangci-lint run ./...
.PHONY: lint

ci: lint test
.PHONY: ci

.DEFAULT_GOAL := ci
