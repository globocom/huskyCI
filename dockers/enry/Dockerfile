# Dockerfile used to create "husyci/enry" image
# https://hub.docker.com/r/huskyci/enry/

FROM golang

RUN apt-get update -y
RUN go get gopkg.in/src-d/enry.v1/...

# installing jq 
RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
RUN chmod +x ./jq
RUN cp jq /usr/bin