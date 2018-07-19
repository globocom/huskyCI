# Dockerfile used to create "husyci/gas" image
# https://hub.docker.com/r/huskyci/gas/

FROM golang

RUN apt-get update -y
RUN go get github.com/GoASTScanner/gas/cmd/gas/...

# installing jq 
RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
RUN chmod +x ./jq
RUN cp jq /usr/bin