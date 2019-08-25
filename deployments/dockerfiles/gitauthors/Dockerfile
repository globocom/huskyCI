# Dockerfile used to create "husyci/gitauthors" image
# https://hub.docker.com/r/huskyci/gitauthors/

FROM alpine:3.8

RUN apk update && apk upgrade \
	&& apk add --update --no-cache git openssh-client
