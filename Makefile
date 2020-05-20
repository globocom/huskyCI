.SILENT:
.DEFAULT_GOAL := help

GO ?= go
GOROOT ?= $(shell $(GO) env GOROOT)
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOCILINT ?= ./bin/golangci-lint
GOLINT ?= $(GOBIN)/golint
GOSEC ?= $(GOBIN)/gosec
GINKGO ?= $(GOBIN)/ginkgo

HUSKYCIBIN ?= huskyci
HUSKYCICLIENTBIN ?= huskyci-client
HUSKYCICLIBIN ?= huskyci-cli

COLOR_RESET = \033[0m
COLOR_COMMAND = \033[36m
COLOR_YELLOW = \033[33m
COLOR_GREEN = \033[32m
COLOR_RED = \033[31m

PROJECT := huskyCI

TAG := $(shell git describe --tags --abbrev=0)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT := $(shell git rev-parse $(TAG))
LDFLAGS := '-X "main.version=$(TAG)" -X "main.commit=$(COMMIT)" -X "main.date=$(DATE)"'

## Builds API code into a binary
build-api:
	$(GO) build -o "$(HUSKYCICLIENTBIN)" /api/server.go

## Builds API code using linux architecture into a binary
build-api-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o "$(HUSKYCICLIENTBIN)" /api/server.go

## Builds client code into a binary
build-client:
	$(GO) build -ldflags $(LDFLAGS) -o "$(HUSKYCIBIN)" client/cmd/main.go

## Builds client code using linux architecture into a binary
build-client-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o "$(HUSKYCICLIENTBIN)" client/cmd/main.go

## Builds cli code into a binary
build-cli:
	$(GO) build -o "$(HUSKYCICLIBIN)" cli/main.go

## Builds cli code using linux architecture into a binary
build-cli-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o "$(HUSKYCICLIBIN)" cli/main.go

## Builds all securityTest containers locally with the latest tags
build-containers:
	chmod +x deployments/scripts/build-containers.sh
	./deployments/scripts/build-containers.sh

## Checks dependencies
check-deps:
	cd api && $(GO) mod tidy && $(GO) mod verify
	cd cli && $(GO) mod tidy && $(GO) mod verify
	cd client && $(GO) mod tidy && $(GO) mod verify

## Runs a security static analysis using Gosec
check-sec: get-gosec-deps gosec

## Checks .env file from huskyCI
check-env:
	cat .env

## Checks every securityTest version from their container images
check-containers-version:
	chmod +x deployments/scripts/check-containers-version.sh
	./deployments/scripts/check-containers-version.sh


## Composes huskyCI environment using docker-compose
compose:
	docker-compose -f deployments/docker-compose.yml down -v
	docker-compose -f deployments/docker-compose.yml up -d --build --force-recreate

## Creates certs and sets all config to huskyCI_Docker_API
create-certs:
	chmod +x deployments/scripts/run-create-certs.sh
	./deployments/scripts/run-create-certs.sh

## Generates a local token to be used in a local environment
generate-local-token:
	chmod +x deployments/scripts/generate-local-token.sh
	./deployments/scripts/generate-local-token.sh

## Generates passwords and set them as environment variables
generate-passwords:
	chmod +x deployments/scripts/generate-env.sh
	./deployments/scripts/generate-env.sh

## Gets all gosec dependencies
get-gosec-deps:
	$(GO) get -u github.com/securego/gosec/cmd/gosec

## Gets all link dependencies
get-lint-deps:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0
	$(GO) get -u golang.org/x/lint/golint

## Runs go lint
golint:
	$(GOLINT) $(shell $(GO) list ./...)

## Runs Golangci-lint
golangci-lint:
	$(GOCILINT) run

## Runs gosec
gosec:
	$(GOSEC) ./... 2> /dev/null

## Prints help message
help:
	printf "\n${COLOR_YELLOW}${PROJECT}\n------\n${COLOR_RESET}"
	awk '/^[a-zA-Z\-\_0-9\.%]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "${COLOR_COMMAND}$$ make %s${COLOR_RESET} %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort
	printf "\n"

## Installs a development environment using docker-compose
install: create-certs prepare-local-mongodb compose generate-passwords generate-local-token

## Runs all huskyCI lint
lint: get-lint-deps golint golangci-lint

## Set up local mongoDB settings file
prepare-local-mongodb:
	chmod +x deployments/scripts/prepare-local-mongodb.sh
	./deployments/scripts/prepare-local-mongodb.sh

## Push securityTest containers to hub.docker
push-containers:
	chmod +x deployments/scripts/push-containers.sh
	./deployments/scripts/push-containers.sh

## Restarts only huskyCI_API container
restart-huskyci-api:
	chmod +x deployments/scripts/restart-huskyci-api.sh
	./deployments/scripts/restart-huskyci-api.sh

## Runs huskyci-client
run-cli: build-cli
	./cli/"$(HUSKYCICLIBIN)" run

## Run huskyci-client compiling it in Linux arch
run-cli-linux: build-cli-linux
	./cli/"$(HUSKYCICLIBIN)" run

## Runs huskyci-client
run-client: build-client
	./client/cmd/"$(HUSKYCICLIENTBIN)"

## Runs huskyci-client with JSON output
run-client-json: build-client
	./client/cmd/"$(HUSKYCICLIENTBIN)" JSON

## Run huskyci-client compiling it in Linux arch
run-client-linux: build-client-linux
	./client/cmd/"$(HUSKYCICLIENTBIN)"

## Run huskyci-client compiling it in Linux arch with JSON output
run-client-linux-json: build-client-linux
	./client/cmd/"$(HUSKYCICLIENTBIN)" JSON

## Performs all unit tests using ginkgo
test:
	cd api && $(GO) test -coverprofile=c.out ./...
	cd api && $(GO) tool cover -func=c.out
	cd api && $(GO) tool cover -html=c.out -o coverage.html

## Builds and push securityTest containers with the latest tags
update-containers: build-containers push-containers
