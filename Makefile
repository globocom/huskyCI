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

## Builds Go project to the executable file huskyci
build:
	cd api && GOOS=linux GOARCH=amd64 $(GO) build -mod vendor -ldflags $(LDFLAGS) -o "$(HUSKYCIBIN)"

## Builds client to the executable file huskyci-client
build-client:
	cd client/cmd && $(GO) build -mod vendor -o "$(HUSKYCICLIENTBIN)" && mv "$(HUSKYCICLIENTBIN)" ../..

## Builds client to the executable file huskyci-client
build-client-linux:
	cd client/cmd && GOOS=linux GOARCH=amd64 $(GO) build -mod vendor -o "$(HUSKYCICLIENTBIN)" && mv "$(HUSKYCICLIENTBIN)" ../..

## Builds CLI to the executable file huskyci-client
build-cli:
	cd cli && $(GO) build -o "$(HUSKYCICLIENTBIN)" main.go

## Builds CLI to the executable file huskyci-client
build-cli-linux:
	cd cli && GOOS=linux GOARCH=amd64 $(GO) build -o "$(HUSKYCICLIENTBIN)" main.go

## Builds all securityTest containers locally with the tag latest
build-containers:
	chmod +x deployments/scripts/build-containers.sh
	./deployments/scripts/build-containers.sh

## Checks depencies of the project
check-deps:
	$(GO) mod verify
	$(GO) mod vendor

## Runs a security static analysis using Gosec
check-sec: get-gosec-deps gosec

## Checks .env file from huskyCI
check-env:
	cat .env

## Checks every securityTest version from their container images
check-containers-version:
	chmod +x deployments/scripts/check-containers-version.sh
	./deployments/scripts/check-containers-version.sh

## Run tests with code coverage
coverage:
	$(GO) test -mod vendor ./... -coverprofile=c.out
	$(GO) tool cover -html=c.out -o coverage.html

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

## Gets all go test dependencies
get-test-deps:
	$(GO) get -u github.com/onsi/ginkgo/ginkgo
	$(GO) get -u github.com/onsi/gomega/...
	$(GO) get -u github.com/mattn/goveralls

## Runs ginkgo
ginkgo:
	$(GINKGO) -r -keepGoing

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
	cd cli && ./"$(HUSKYCICLIENTBIN)" run

## Run huskyci-client compiling it in Linux arch
run-cli-linux: build-cli-linux
	cd cli && ./"$(HUSKYCICLIENTBIN)" run

## Runs huskyci-client
run-client: build-client
	./"$(HUSKYCICLIENTBIN)"

## Runs huskyci-client with JSON output
run-client-json: build-client
	./"$(HUSKYCICLIENTBIN)" JSON

## Run huskyci-client compiling it in Linux arch
run-client-linux: build-client-linux
	./"$(HUSKYCICLIENTBIN)"

## Run huskyci-client compiling it in Linux arch with JSON output
run-client-linux-json: build-client-linux
	./"$(HUSKYCICLIENTBIN)" JSON

## Perfoms all make tests
test: get-test-deps ginkgo coverage

## Builds and push securityTest containers with the latest tags
update-containers: build-containers push-containers
