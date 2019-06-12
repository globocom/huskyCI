# Dockerfile used to create "husyci/gosec" image
# https://hub.docker.com/r/huskyci/gosec/

FROM golang:alpine 

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk git \
	&& apk add --update --no-cache git openssh-client jq

RUN go get github.com/securego/gosec/cmd/gosec/...

