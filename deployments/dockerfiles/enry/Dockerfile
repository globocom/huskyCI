# Dockerfile used to create "husyci/enry" image
# https://hub.docker.com/r/huskyci/enry/
FROM golang:alpine as builder

RUN apk update && apk upgrade \
	&& apk add --update --no-cache git openssh-client

RUN git clone https://github.com/go-enry/enry.git && \
    cd enry && \
	go mod tidy && \
	CGO_ENABLED=0 go build && \
	cp enry /usr/bin

FROM alpine:latest

RUN apk update && apk upgrade \
	&& apk add --update --no-cache git openssh-client

COPY --from=builder /usr/bin/enry /usr/bin
