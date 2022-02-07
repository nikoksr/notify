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

fmt:
	@go install mvdan.cc/gofumpt@latest
	@gofumpt -w -l .

	@go install github.com/daixiang0/gci@latest
	@gci write --Section Standard --Section Default --Section "Prefix(github.com/nikoksr/notify)" .
.PHONY: fmt


lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run --config .golangci.yml
.PHONY: lint

ci: lint test
.PHONY: ci

###############################################################################
# DEPENDENCIES
###############################################################################

changelog-latest:
	@git cliff --config .cliff.toml --latest --strip header
.PHONY: changelog-latest

changelog-file:
	@git cliff --config .cliff.toml --topo-order -o CHANGELOG.md
.PHONY: changelog-file

###############################################################################

.DEFAULT_GOAL := ci
