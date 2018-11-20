# Dockerfile used to create "husyci/gas" image
# https://hub.docker.com/r/huskyci/gas/

FROM golang:alpine

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk bash openssh-client \
	&& apk add git

RUN go get github.com/securego/gosec/cmd/gosec/...

RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
RUN chmod +x ./jq
RUN cp jq /usr/bin
