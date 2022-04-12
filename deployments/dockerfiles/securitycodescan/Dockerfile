FROM mcr.microsoft.com/dotnet/sdk:6.0-alpine

RUN apk add jq openssh bash

RUN dotnet tool install --global security-scan --version 5.6.2

RUN wget https://dot.net/v1/dotnet-install.sh
RUN chmod +x dotnet-install.sh
RUN ./dotnet-install.sh -c 3.1 --install-dir /usr/share/dotnet
RUN ./dotnet-install.sh -c 5.0 --install-dir /usr/share/dotnet
RUN ./dotnet-install.sh -c 7.0 --install-dir /usr/share/dotnet

ENV PATH="$PATH:/root/.dotnet/tools"