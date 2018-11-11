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

## Installs all development dependencies
install-deps:
	$(GO) get -u github.com/golang/dep/cmd/dep
	$(GO) get -u golang.org/x/lint/golint
	$(GO) get -u github.com/securego/gosec/cmd/gosec

## Runs a security static analysis using Gosec
security-check:
	$(GOSEC) ./... 2> /dev/null

## Perfoms all make tests
test: lint security-check

## Runs lint
lint:
	$(GOLINT) $(shell $(GO) list ./...)

## Runs project
run:
	@go run server.go

## Builds Go project to the executable file huskyci
build:
	$(GO) build -o "$(HUSKYCIBIN)"

## Builds a development environment using docker-compose
compose-start:
	docker-compose -f dev/docker-compose.yml build
	docker-compose -f dev/docker-compose.yml up -d

## Stops all dockers using docker-compose
compose-stop:
	docker-compose -f dev/docker-compose.yml stop

# Restarts all dockers from the development environment
compose-restart: compose-stop compose-start

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
