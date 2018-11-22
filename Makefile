.SILENT:
.DEFAULT_GOAL := help

GO ?= go
GOROOT ?= $(shell $(GO) env GOROOT)
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin
GODEP ?= $(GOBIN)/dep
GOLINT ?= $(GOBIN)/golint
GOSEC ?= $(GOBIN)/gosec

HUSKYCIBIN ?= huskyci

COLOR_RESET = \033[0m
COLOR_COMMAND = \033[36m
COLOR_YELLOW = \033[33m
COLOR_GREEN = \033[32m
COLOR_RED = \033[31m

PROJECT := HuskyCI

TAG := $(shell git describe --tags --abbrev=0)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT := $(shell git rev-parse $(TAG))
LDFLAGS := '-X "main.version=$(TAG)" -X "main.commit=$(COMMIT)" -X "main.date=$(DATE)"'

## Installs a development environment using docker-compose
install: generate-passwords create-certs compose check-env pull-images

## Gets all go test dependencies
get-deps:
	$(GO) get -u github.com/golang/dep/cmd/dep
	$(GO) get -u golang.org/x/lint/golint
	$(GO) get -u github.com/securego/gosec/cmd/gosec

## Checks depencies of the project
check-deps:
	$(GODEP) ensure -v

## Runs a security static analysis using Gosec
check-sec:
	$(GOSEC) ./... 2> /dev/null

## Checks .env file from HuskyCI
check-env:
	cat .env

## Perfoms all make tests
test: get-deps lint security-check

## Runs lint
lint:
	$(GOLINT) $(shell $(GO) list ./...)

## Builds Go project to the executable file huskyci
build:
	$(GO) build -ldflags $(LDFLAGS) -o "$(HUSKYCIBIN)"


## Run project using docker-compose
compose:
	docker-compose -f dev/docker-compose.yml build
	docker-compose -f dev/docker-compose.yml down
	docker-compose -f dev/docker-compose.yml up -d --force-recreate
	docker-compose -f dev/docker-compose.yml logs -f huskyapi

## Pulls every HuskyCI docker image into dockerAPI container
pull-images:
	docker exec dockerAPI /bin/sh -c "docker pull huskyci/enry"
	docker exec dockerAPI /bin/sh -c "docker pull huskyci/gas"
	docker exec dockerAPI /bin/sh -c "docker pull huskyci/bandit"
	docker exec dockerAPI /bin/sh -c "docker pull huskyci/brakeman"
	#docker exec dockerAPI /bin/sh -c "docker pull huskyci/retirejs"

## Creates certs and sets all config to dockerAPI
create-certs:
	chmod +x scripts/run-create-certs.sh
	./scripts/run-create-certs.sh

## Generates passwords and set them as environment variables
generate-passwords:
	chmod +x scripts/generate-env-pass.sh
	./scripts/generate-env-pass.sh

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
