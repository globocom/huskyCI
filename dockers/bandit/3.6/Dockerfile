FROM docker-hub.artifactory.globoi.com/python:3.6-alpine


RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk bash openssh-client \
	&& apk add git

RUN pip3 install bandit