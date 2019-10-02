# Dockerfile used to create "husyci/safety" image
# https://hub.docker.com/r/huskyci/safety/

FROM python:3.6-alpine

COPY ./script.sh /
RUN chmod +x /script.sh

RUN apk add --virtual .build-deps --no-cache git alpine-sdk bash openssh-client \
    && pip3 install safety \
    && apk del .build-deps # remove virtual namespace after building is done 

# Installing https://github.com/stedolan/jq
RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 \
    && chmod +x ./jq && cp jq /usr/bin
