# ==============================================================================
# Install dependencies

dev.setup.mac:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli

# ==============================================================================
# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := 1.0

# ==============================================================================
# Running tests within the local computer
# go install honnef.co/go/tools/cmd/staticcheck@latest
# go install golang.org/x/vuln/cmd/govulncheck@latest

test:
	go test -count=1 ./...
	staticcheck -checks=all ./...
	govulncheck ./...

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Docker support

docker-down:
	docker rm -f $(shell docker ps -aq)

docker-clean:
	docker system prune -f	

docker-kind-logs:
	docker logs -f $(KIND_CLUSTER)-control-plane

# ==============================================================================
# Dev support

dev:
	docker-compose -f zarf/dev/docker-compose.yaml up -d --build

dev-stop:
	docker-compose -f zarf/dev/docker-compose.yaml down