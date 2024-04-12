DATABASE_DSN = postgresql://localhost:5432/postgres

DOC_PORT = 6060

SERVER_BINARY_NAME = server
SERVER_PACKAGE_PATH = ./cmd/server
SERVER_ADDR = localhost:8080

DEFAULT_EXPIRATION = 5m
CLEANUP_INTERVAL = 10m

# ====================
# HELPERS
# ====================

## help: show this help message
help:
	@echo
	@echo 'usage: make <target>'
	@echo
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
	@echo

# ====================
# QUALITY
# ====================

## tidy: format code and tidy mod file
tidy:
	go fmt ./...
	go mod tidy -v

# ====================
# DEVELOPMENT
# ====================

## build-local: build the server locally
build-local:
	go build -o /tmp/bin/$(SERVER_BINARY_NAME) $(SERVER_PACKAGE_PATH)

## run-local: run the server locally
run-local: build-local
	/tmp/bin/$(SERVER_BINARY_NAME) -a=$(SERVER_ADDR) -d=$(DATABASE_DSN) -clean=$(CLEANUP_INTERVAL) -exp=$(DEFAULT_EXPIRATION)

## build-docker: build the server with docker-compose
build-docker:
	docker-compose build

## run-docker: build the server with docker-compose
run-docker: build-docker
	docker-compose up

# ====================
# DOCUMENTATION
# ====================

## doc: generate documentation on http port
doc:
	@echo 'open http://localhost:$(DOC_PORT)/pkg/github.com/pavlegich/banners-service/?m=all'
	godoc -http=:$(DOC_PORT)

.PHONY: help tidy build-local run-local build-docker run-docker doc
