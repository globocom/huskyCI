# Dockerfile used to create "husyci/bandit:3.6" image
# https://hub.docker.com/r/huskyci/bandit/

FROM python:3.6-alpine

# Installing bandit and jq
RUN apk add --no-cache git jq alpine-sdk bash openssh-client \ 
    && pip3 install bandit
