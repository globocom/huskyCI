# Dockerfile used to create "husyci/gosec" image
# https://hub.docker.com/r/huskyci/gosec/

FROM golang:alpine as builder

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk git \
	&& apk --update add jq

RUN go get github.com/securego/gosec/cmd/gosec/...

# From the base image
FROM alpine:3.8
COPY --from=builder /go/bin/gosec /usr/bin/
COPY --from=builder /go/jq /usr/bin/
