# Dockerfile used to create "husyci/gosec" image
# https://hub.docker.com/r/huskyci/gosec/

FROM golang:alpine as builder

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk git \
	&& apk add --update --no-cache git openssh-client jq

RUN go get github.com/securego/gosec/cmd/gosec/...

FROM alpine:latest

RUN apk update && apk upgrade && apk --no-cache add ca-certificates git openssh-client jq

WORKDIR /go

COPY --from=builder /go/bin/gosec /go/bin/gosec
COPY --from=builder /usr/local/go/bin/go /usr/local/go/bin/go
COPY --from=builder /usr/local/go/bin/gofmt /usr/local/go/bin/gofmt
COPY --from=builder /usr/local/go/src/ /usr/local/go/src/

ENV PATH=/go/bin:/usr/local/go/bin:$PATH
