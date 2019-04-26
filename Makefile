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
	cd api && GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o "$(HUSKYCIBIN)"

## Builds client to the executable file huskyci-client
build-client:
	cd client/cmd && $(GO) build -o "$(HUSKYCICLIENTBIN)" && mv "$(HUSKYCICLIENTBIN)" ../..

## Builds client to the executable file huskyci-client
build-client-linux:
	cd client/cmd && GOOS=linux GOARCH=amd64 $(GO) build -o "$(HUSKYCICLIENTBIN)" && mv "$(HUSKYCICLIENTBIN)" ../..

## Checks depencies of the project
check-deps:
	$(GODEP) ensure -v

## Runs a security static analysis using Gosec
check-sec:
	$(GOSEC) ./... 2> /dev/null

## Checks .env file from huskyCI
check-env:
	cat .env

## Composes huskyCI environment using docker-compose
compose:
	docker-compose -f deployments/docker-compose.yml down -v
	docker-compose -f deployments/docker-compose.yml up -d --build --force-recreate

## Creates certs and sets all config to huskyCI_Docker_API
create-certs:
	chmod +x deployments/scripts/run-create-certs.sh
	./deployments/scripts/run-create-certs.sh

## Generates passwords and set them as environment variables
generate-passwords:
	chmod +x deployments/scripts/generate-env-pass.sh
	./deployments/scripts/generate-env-pass.sh

## Gets all go test dependencies
get-test-deps:
	$(GO) get -u github.com/golang/dep/cmd/dep
	$(GO) get -u golang.org/x/lint/golint
	$(GO) get -u github.com/securego/gosec/cmd/gosec
	$(GO) get -u github.com/onsi/ginkgo/ginkgo
	$(GO) get -u github.com/onsi/gomega/...

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
install: generate-passwords create-certs compose

## Installs a development environment using docker-compose but does not pull security tests' images
install-pull-images: generate-passwords create-certs compose pull-images

## Runs lint
lint:
	$(GOLINT) $(shell $(GO) list ./...)

## Pulls every HuskyCI docker image into huskyCI_Docker_API container
pull-images:
	docker exec huskyCI_Docker_API /bin/sh -c "docker pull huskyci/enry"
	docker exec huskyCI_Docker_API /bin/sh -c "docker pull huskyci/gosec"
	docker exec huskyCI_Docker_API /bin/sh -c "docker pull huskyci/bandit"
	docker exec huskyCI_Docker_API /bin/sh -c "docker pull huskyci/brakeman"
	docker exec huskyCI_Docker_API /bin/sh -c "docker pull huskyci/retirejs"
	docker exec huskyCI_Docker_API /bin/sh -c "docker pull huskyci/safety"

## Runs huskyci-client
run-client: build-client
	./"$(HUSKYCICLIENTBIN)"

## Runs ginkgo
ginkgo:
	ginkgo -r -keepGoing 
	
## Perfoms all make tests
test: get-test-deps lint check-sec ginkgo
