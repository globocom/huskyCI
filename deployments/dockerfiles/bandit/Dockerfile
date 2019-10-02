# Dockerfile used to create "husyci/bandit:3.6" image
# https://hub.docker.com/r/huskyci/bandit/

FROM python:3.6-alpine

# Installing bandit via pip
RUN apk add --virtual .build-deps --no-cache git alpine-sdk bash openssh-client \
    && pip3 install bandit \
    && apk del .build-deps # remove virtual namespace after building is done 

# Installing https://github.com/stedolan/jq
RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 \
    && chmod +x ./jq && cp jq /usr/bin
