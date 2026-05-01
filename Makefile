BINARY := mdview
BUILD_DIR := bin
MAIN_PKG := .
GO_PACKAGES := ./...

.PHONY: help lint fmt fmt-check test build run clean check

help:
	@printf '%s\n' 'Available targets:'
	@printf '  %-10s %s\n' 'lint' 'Run go vet'
	@printf '  %-10s %s\n' 'fmt' 'Format Go files with gofmt'
	@printf '  %-10s %s\n' 'fmt-check' 'Check Go formatting without modifying files'
	@printf '  %-10s %s\n' 'test' 'Run all tests'
	@printf '  %-10s %s\n' 'build' 'Build the mdview binary into bin/'
	@printf '  %-10s %s\n' 'run' 'Run mdview; pass FILE=path/to/file.md'
	@printf '  %-10s %s\n' 'check' 'Run fmt-check, lint, and test'
	@printf '  %-10s %s\n' 'clean' 'Remove build artifacts'

lint:
	go vet $(GO_PACKAGES)

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

fmt-check:
	@test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*'))"

test:
	go test $(GO_PACKAGES)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) $(MAIN_PKG)

run:
	@test -n "$(FILE)" || (printf '%s\n' 'usage: make run FILE=path/to/file.md' >&2; exit 2)
	go run $(MAIN_PKG) "$(FILE)"

check: fmt-check lint test

clean:
	rm -rf $(BUILD_DIR)
