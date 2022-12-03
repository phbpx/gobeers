
SHELL := /bin/bash
ARGS = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`
VERSION := 1.0

# ==============================================================================
# help

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# ==============================================================================
# Setup

GOBIN := $(shell go version)

check.go:
	@go version >/dev/null 2>&1 || (echo "ERROR: go is not installed" && exit 1)

check.docker:
	@docker version >/dev/null 2>&1 || (echo "ERROR: docker is not installed" && exit 1)
	@docker-compose version >/dev/null 2>&1 || (echo "ERROR: docker-compose is not installed" && exit 1)

## Install go tools
setup: check.go check.docker
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@docker pull postgres:15-alpine

# ==============================================================================
# Running tests within the local computer
# go install honnef.co/go/tools/cmd/staticcheck@latest
# go install golang.org/x/vuln/cmd/govulncheck@latest

## Run tests in the local computer
test:
	go test -count=1 ./...
	staticcheck -checks=all ./...
	govulncheck ./...

# ==============================================================================
# Modules support

## Install and vendor project dependencies
tidy:
	go mod tidy
	go mod vendor

## Reset project dependencies
deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

## Upgrade project dependencies
deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Administration

## Create a new database migration into: business/data/dbscheme/sql
add-migration:
	migrate create -dir ./business/data/dbschema/sql -ext sql -seq $(call args,defaultstring)

# ==============================================================================
# Docker support

docker-down:
	docker rm -f $(shell docker ps -aq)

docker-clean:
	docker system prune -f	

# ==============================================================================
# Dev support

dev:
	docker-compose -f zarf/dev/docker-compose.yaml up -d --build

dev-stop:
	docker-compose -f zarf/dev/docker-compose.yaml down