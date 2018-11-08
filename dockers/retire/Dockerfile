FROM docker-hub.artifactory.globoi.com/node:alpine

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk bash openssh-client \
	&& apk add git


RUN npm install -g retire


RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
RUN chmod +x ./jq
RUN cp jq /usr/bin