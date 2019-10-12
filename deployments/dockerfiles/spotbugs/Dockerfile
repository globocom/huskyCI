FROM openjdk:8-jdk-alpine

RUN apk add --no-cache curl tar bash wget unzip

ARG MAVEN_VERSION=3.6.2
ARG USER_HOME_DIR="/root"
ARG SHA=3fbc92d1961482d6fbd57fbf3dd6d27a4de70778528ee3fb44aa7d27eb32dfdc
ARG BASE_URL=https://apache.osuosl.org/maven/maven-3/${MAVEN_VERSION}/binaries

RUN mkdir -p /usr/share/maven /usr/share/maven/ref \
  && curl -fsSL -o /tmp/apache-maven.tar.gz ${BASE_URL}/apache-maven-${MAVEN_VERSION}-bin.tar.gz \
  && echo "${SHA}  /tmp/apache-maven.tar.gz" | sha256sum -c - \
  && tar -xzf /tmp/apache-maven.tar.gz -C /usr/share/maven --strip-components=1 \
  && rm -f /tmp/apache-maven.tar.gz \
  && ln -s /usr/share/maven/bin/mvn /usr/bin/mvn \
  && mkdir -p /opt \
  && cd /opt \
  && export SPOTBUGS_VERSION=$(curl --silent "https://api.github.com/repos/spotbugs/spotbugs/releases/latest"|grep '"tag_name":'|sed -E 's/.*"([^"]+)".*/\1/'|cut -d '/' -f2|sed  's/_/-/') \
  && wget -nc -O spotbugs.zip http://repo.maven.apache.org/maven2/com/github/spotbugs/spotbugs/${SPOTBUGS_VERSION}/spotbugs-${SPOTBUGS_VERSION}.zip \
  && unzip spotbugs.zip \
  && rm -f spotbugs.zip \
  && mv spotbugs-${SPOTBUGS_VERSION} spotbugs \
  && wget -nc -O findsecbugs-plugin-1.9.0.jar https://repo1.maven.org/maven2/com/h3xstream/findsecbugs/findsecbugs-plugin/1.9.0/findsecbugs-plugin-1.9.0.jar \
  && echo $SPOTBUGS_VERSION > /opt/spotbugs/version

ENV MAVEN_HOME /usr/share/maven
ENV MAVEN_CONFIG "$USER_HOME_DIR/.m2"

COPY mvn-entrypoint.sh /usr/local/bin/mvn-entrypoint.sh
COPY settings-docker.xml /usr/share/maven/ref/

VOLUME "$USER_HOME_DIR/.m2"

