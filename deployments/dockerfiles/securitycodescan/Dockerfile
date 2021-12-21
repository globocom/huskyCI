FROM mcr.microsoft.com/dotnet/sdk:6.0-alpine

RUN apk add jq openssh

RUN dotnet tool install --global security-scan --version 5.6.0

ENV PATH="$PATH:/root/.dotnet/tools"