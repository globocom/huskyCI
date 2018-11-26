# Dockerfile used to create "husyci/enry" image
# https://hub.docker.com/r/huskyci/enry/

FROM golang:alpine as builder

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk git

RUN go get gopkg.in/src-d/enry.v1/...

RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
RUN chmod +x ./jq

# From the base image
FROM alpine:latest
COPY --from=builder /go/bin/enry /usr/bin/
COPY --from=builder /go/jq /usr/bin/