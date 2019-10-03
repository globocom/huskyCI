# Dockerfile used to create "husyci/gosec" image
# https://hub.docker.com/r/huskyci/gosec/

FROM securego/gosec:latest

RUN apk update && apk upgrade && apk --no-cache add ca-certificates openssh-client jq

ENTRYPOINT [ "/bin/sh" ]
