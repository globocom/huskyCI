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

## Installs a development environment using docker-compose
install: compose-start pull-images config-auth

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

## Perfoms all make tests
test: get-deps lint security-check

## Runs lint
lint:
	$(GOLINT) $(shell $(GO) list ./...)

## Builds Go project to the executable file huskyci
build:
	$(GO) build -o "$(HUSKYCIBIN)"

## Builds an environment using docker-compose
compose-start:
	docker-compose -f dev/docker-compose.yml build
	docker-compose -f dev/docker-compose.yml up -d

## Stops all dockers using docker-compose
compose-stop:
	docker-compose -f dev/docker-compose.yml stop

## Restarts all dockers from the development environment
compose-restart: compose-stop compose-start

## Destroys all dockers from the development environment
compose-destroy:
	docker stop `docker ps -a | grep 'huskyCIAPI\|mongodb\|dockerAPI' | awk '{print $$1}'`
	docker rm `docker ps -a | grep 'huskyCIAPI\|mongodb\|dockerAPI' | awk '{print $$1}'`

## Pulls every HuskyCI docker image into dockerAPI container
pull-images:
	docker exec dockerAPI /bin/sh -c "docker pull huskyci/enry"
	docker exec dockerAPI /bin/sh -c "docker pull huskyci/gas"
	docker exec dockerAPI /bin/sh -c "docker pull huskyci/bandit:2.7"
	docker exec dockerAPI /bin/sh -c "docker pull huskyci/bandit:3.6"
	docker exec dockerAPI /bin/sh -c "docker pull huskyci/brakeman"
	docker exec dockerAPI /bin/sh -c "docker pull huskyci/retirejs"

## Configures DockerAPI so that it only accepts requests from HuskyCI
setup-auth:
	echo "[DockerAPI - Server]"
	openssl genrsa -aes256 -out ca-key.pem 4096
	echo "Step 2"
	openssl req -new -x509 -days 365 -key ca-key.pem -sha256 -out ca.pem -subj "/C=US/ST=New Sweden/L=Stockholm /O=.../OU=.../CN=.../emailAddress=..."
	echo "Step 3"
	openssl genrsa -out server-key.pem 4096
	openssl req -subj "/CN=dockerAPI" -sha256 -new -key server-key.pem -out server.csr
	echo subjectAltName = DNS:dockerAPI >> extfile.cnf
	echo extendedKeyUsage = serverAuth >> extfile.cnf
	openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile extfile.cnf

	echo "[HuskyCI - Client]"
	openssl genrsa -out key.pem 4096
	echo "Step 7"
	openssl req -subj '/CN=client' -new -key key.pem -out client.csr
	echo "Step 8"
	echo extendedKeyUsage = clientAuth >> extfile2.cnf
	echo "Step 9"
	openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out cert.pem -extfile extfile2.cnf
	echo "Step 10"
	chmod 0400 ca-key.pem key.pem server-key.pem
	chmod 0444 ca.pem server-cert.pem cert.pem

	echo "[Sending pem files to container dockerAPI]"
	docker cp ca-key.pem dockerAPI:/ca-key.pem
	docker cp ca.pem dockerAPI:/ca.pem
	docker cp ca.srl dockerAPI:/ca.srl
	docker cp server-key.pem dockerAPI:/server-key.pem
	docker cp server-cert.pem dockerAPI:/server-cert.pem
	docker exec dockerAPI /bin/sh -c "chown root:root /*.pem"

	echo "[Sending pem files to container huskyCIAPI]"
	docker cp key.pem huskyCIAPI:/go/src/github.com/globocom/husky/key.pem
	docker cp cert.pem huskyCIAPI:/go/src/github.com/globocom/husky/cert.pem
	docker cp ca.pem huskyCIAPI:/go/src/github.com/globocom/husky/ca.pem
	docker exec huskyCIAPI /bin/sh -c "chown root:root /go/src/github.com/globocom/husky/*.pem"

	echo "[Cleaning]"
	rm -rf client.csr server.csr extfile.cnf extfile2.cnf ca.srl *.pem

	# [dockerAPI] # dockerd --tlsverify --tlscacert=ca.pem --tlscert=server-cert.pem --tlskey=server-key.pem -H=0.0.0.0:2375
	# Error starting daemon: pid file found, ensure docker is not running or delete /var/run/docker.pid

## Run project using docker-compose
compose:
	docker-compose -f dev/docker-compose.yml build
	docker-compose -f dev/docker-compose.yml down -v
	docker-compose -f dev/docker-compose.yml up -d --force-recreate


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
