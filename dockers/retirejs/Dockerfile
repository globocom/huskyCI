# Dockerfile used to create "huskyci/retirejs" image
# https://hub.docker.com/r/huskyci/retirejs/

FROM node:alpine

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk bash openssh-client \
	&& apk add git

RUN npm install
RUN npm install -g retire
