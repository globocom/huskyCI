FROM golang

ADD . /go/src/github.com/globocom/huskyci/api
WORKDIR /go/src/github.com/globocom/huskyci/api
