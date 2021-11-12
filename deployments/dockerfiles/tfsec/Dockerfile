# Dockerfile used to create "husyci/tfsec" image
# https://hub.docker.com/r/huskyci/tfsec/
FROM tfsec/tfsec:v0.58

USER root

RUN apk update && apk upgrade \
	&& apk add git jq openssh-client curl