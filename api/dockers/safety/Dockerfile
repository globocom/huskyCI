# Dockerfile used to create "husyci/safety" image
# https://hub.docker.com/r/huskyci/safety/

FROM python:3.6-alpine

COPY ./script.sh /

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk bash openssh-client \
	&& apk add git

RUN pip3 install safety

RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
RUN chmod +x ./jq
RUN cp jq /usr/bin