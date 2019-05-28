RM_F := rm -f
GO := go
LINT := $(GOPATH)/bin/golangci-lint

Additional_Linters := misspell,nakedret

export GO111MODULE=on

.PHONY: help
help: ## Show this help message
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
        	printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
        }' $(MAKEFILE_LIST)

build: ## Build
	$(GO) build -o githubfs ./cmd

.PHONY: tidy
tidy: ## Run linters
	$(GO) mod tidy
	$(LINT) --enable $(Additional_Linters) run

.PHONY: tools
tools: ## Install build tools
	$(GO) get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0

.PHONY: check
check: ## Run tests
	$(GO) test -v ./cmd

.PHONY: clean
clean: ## Remove output files
	$(RM_F) githubfs

