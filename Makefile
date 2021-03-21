export GO111MODULE := on
export GOPROXY = https://proxy.golang.org,direct

###############################################################################
# DEPENDENCIES
###############################################################################

# Install all the build and lint dependencies
setup:
	go mod download
	go generate -v ./...
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

# gofumports and gci all go files
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofumpt -w "$$file"; done
	gci -w -local github.com/nikoksr/notify .
.PHONY: fmt

# Run all the linters
lint:
	golangci-lint run ./...
.PHONY: lint

ci: test lint
.PHONY: ci

.DEFAULT_GOAL := ci
