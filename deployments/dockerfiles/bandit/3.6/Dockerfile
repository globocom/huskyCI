# Dockerfile used to create "husyci/bandit:3.6" image
# https://hub.docker.com/r/huskyci/bandit/

FROM python:3.6-alpine

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk bash openssh-client \
	&& apk add git

RUN pip3 install bandit

RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
RUN chmod +x ./jq
RUN cp jq /usr/bin
