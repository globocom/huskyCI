FROM golang:alpine as builder

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk git \
	&& apk add --update --no-cache git openssh-client jq

RUN go get github.com/securego/gosec/cmd/gosec/...

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /go

COPY --from=builder /go/bin/gosec /go/bin/gosec
ENV PATH=/go/bin:$PATH
