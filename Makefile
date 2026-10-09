.DEFAULT_GOAL := help

STATICCHECK_VERSION := 2026.2.1

.PHONY: help
help: ## Outputs the help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Runs all unit tests
	go test -v -race ./...

.PHONY: vet
vet: ## Runs go vet
	go vet ./...

.PHONY: staticcheck
staticcheck: ## Runs static code analyzer staticcheck
	go install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)
	staticcheck ./...