# Dockerfile used to create "husyci/tfsec" image
# https://hub.docker.com/r/huskyci/tfsec/
FROM golang:1.13-alpine

ARG TFSEC_VERSION=v0.19.0

RUN apk update && apk upgrade \
	&& apk add git jq openssh-client

RUN wget https://github.com/liamg/tfsec/releases/download/${TFSEC_VERSION}/tfsec-linux-amd64

RUN mv tfsec-linux-amd64 tfsec

RUN chmod +x tfsec