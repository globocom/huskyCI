# Dockerfile used to create "husyci/tfsec" image
# https://hub.docker.com/r/huskyci/tfsec/
FROM golang:1.13-alpine

RUN apk update && apk upgrade \
	&& apk add git jq openssh-client curl

RUN set -o pipefail && curl https://api.github.com/repos/liamg/tfsec/releases/latest | jq -r ".assets[] | select(.name | contains(\"tfsec-linux-amd64\")) | .browser_download_url" | xargs wget

RUN mv tfsec-linux-amd64 tfsec

RUN chmod +x tfsec