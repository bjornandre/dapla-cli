.PHONY: default
default: | help

.PHONY: run
run: ## Run dapla-cli (without build)
	@go run main.go

.PHONY: changelog
changelog: ## Generate CHANGELOG.md
	github_changelog_generator -u statisticsnorway -p dapla-cli

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
