# Dockerfile used to create "huskyci/npmaudit" image
# https://hub.docker.com/r/huskyci/infer/
#
# Dockerfile highly inspired by:
# - https://github.com/facebook/infer/tree/main/docker
# - https://github.com/globocom/huskyCI/tree/master/deployments/dockerfiles/spotbugs
FROM debian:bullseye-slim

ARG MAVEN_VERSION=3.8.3
ARG GRADLE_VERSION=5.6.2

# mkdir the man/man1 directory due to Debian bug #863199
RUN apt-get update \
    && mkdir -p /usr/share/man/man1 \
    && apt-get install --yes --no-install-recommends \
      curl \
      wget \
      unzip \
      libc6-dev \
      openjdk-11-jdk-headless \
      sqlite3 \
      xz-utils \
      zlib1g-dev \
    && mkdir -p /usr/share/maven /usr/share/maven/ref \
    && curl -fsSL -o /tmp/apache-maven.tar.gz https://apache.osuosl.org/maven/maven-3/${MAVEN_VERSION}/binaries/apache-maven-${MAVEN_VERSION}-bin.tar.gz \
    && tar -xzf /tmp/apache-maven.tar.gz -C /usr/share/maven --strip-components=1 \
    && rm -f /tmp/apache-maven.tar.gz \
    && ln -s /usr/share/maven/bin/mvn /usr/bin/mvn \
    && mkdir -p /opt \
    && cd /opt \
    && wget -nc -O gradle.zip https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip \
    && unzip gradle.zip \
    && rm -f gradle.zip \
    && mv gradle-${GRADLE_VERSION} gradle \
    && rm -rf /var/lib/apt/lists/*

COPY mvn-entrypoint.sh /usr/local/bin/mvn-entrypoint.sh

# Download the Infer release
RUN INFER_VERSION=v1.1.0; \
    cd /opt \
    && curl -sL \
      https://github.com/facebook/infer/releases/download/${INFER_VERSION}/infer-linux64-${INFER_VERSION}.tar.xz | \
    tar xJ \
    && rm -f /infer \
    && ln -s ${PWD}/infer-linux64-$INFER_VERSION /infer

# Install infer
ENV PATH /infer/bin:${PATH}
