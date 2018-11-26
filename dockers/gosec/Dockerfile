# Dockerfile used to create "husyci/gosec" image
# https://hub.docker.com/r/huskyci/gosec/

FROM golang:alpine as builder

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk git

RUN go get github.com/securego/gosec/cmd/gosec/...

RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
RUN chmod +x ./jq

# From the base image
FROM alpine:latest
COPY --from=builder /go/bin/gosec /usr/bin/
COPY --from=builder /go/jq /usr/bin/
