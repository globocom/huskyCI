FROM docker-hub.artifactory.globoi.com/python:2.7-alpine

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk bash openssh-client \
	&& apk add git

RUN pip install bandit
