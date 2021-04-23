.PHONY: default
default: | help

.PHONY: build-local
version = "localbuild"
sha1hash := $(shell git rev-parse --short HEAD)
build_time := $(shell date -u +%Y-%m-%dT%H:%M:%S%Z)
build-local: ## Build dapla-cli
	go build \
	--ldflags "-X github.com/statisticsnorway/dapla-cli/cmd.Version=$(version) \
	-X github.com/statisticsnorway/dapla-cli/cmd.GitSha1Hash=$(sha1hash) \
	-X github.com/statisticsnorway/dapla-cli/cmd.BuildTime=$(build_time)" \
	-o bin/dapla-cli .

.PHONY: build-docker
build-docker: ## Build dapla-cli with docker
	docker build -t dapla-cli .

.PHONY: changelog
changelog: ## Generate CHANGELOG.md
	github_changelog_generator -u statisticsnorway -p dapla-cli

.PHONY: alias-dev
alias-dev: ## Print dapla alias for local development (no build) - apply with eval $(make alias-dev)
	@echo "alias dapla=\"go run main.go --config .dapla-cli-localdev.yml\""

.PHONY: alias-localbuild
alias-localbuild: ## Print dapla alias for local build - apply with eval $(make alias-localbuild)
	@echo "alias dapla=\"bin/dapla-cli\""

.PHONY: alias-docker
alias-docker: ## Print dapla alias for running dapla-cli within docker - apply with eval $(make alias-docker)
	@echo "alias dapla=\"docker run -p 3000:3000 dapla-cli\""

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
