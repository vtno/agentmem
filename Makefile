.DEFAULT_GOAL := help
BINARY        := agentmem-site
VERSION       ?= dev
LDFLAGS       := -ldflags "-s -w -X main.version=$(VERSION)"
INSTALL_DIR   ?= $(HOME)/.local/bin

.PHONY: build test fmt clean install run tidy help

build: ## build the agentmem-site binary (site/ embedded)
	go build -trimpath $(LDFLAGS) -o $(BINARY) .

test: ## run all tests with race detector
	go test -race ./...

fmt: ## format all Go source files
	gofmt -w .

clean: ## remove build artifacts
	rm -f $(BINARY)

install: ## install binary to $(INSTALL_DIR) (VERSION=x.y.z for a tagged build)
	@mkdir -p $(INSTALL_DIR)
	go build -trimpath $(LDFLAGS) -o $(INSTALL_DIR)/$(BINARY) .
	@echo "installed $(INSTALL_DIR)/$(BINARY)"

run: ## run the site locally (default :5555; override ADDR=:8080)
	go run . -addr $(or $(ADDR),:5555)

tidy: ## tidy go.mod and go.sum
	go mod tidy

help: ## list available targets
	@grep -E '^[a-zA-Z_/-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'
