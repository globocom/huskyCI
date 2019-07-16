# Dockerfile used to create "huskyci/npmaudit" image
# https://hub.docker.com/r/huskyci/npmaudit/

FROM node:alpine

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk bash openssh-client \
	&& apk add git