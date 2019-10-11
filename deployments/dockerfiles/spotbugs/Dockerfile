FROM openjdk:13-alpine

RUN apk add --update --no-cache wget unzip curl
RUN mkdir -p /opt
RUN cd /opt \
      && export SPOTBUGS_VERSION=$(curl --silent "https://api.github.com/repos/spotbugs/spotbugs/releases/latest"|grep '"tag_name":'|sed -E 's/.*"([^"]+)".*/\1/'|cut -d '/' -f2|sed  's/_/-/');wget -nc -O spotbugs.zip http://repo.maven.apache.org/maven2/com/github/spotbugs/spotbugs/${SPOTBUGS_VERSION}/spotbugs-${SPOTBUGS_VERSION}.zip \
      && unzip spotbugs.zip \
      && rm -f spotbugs.zip \
      && mv spotbugs-${SPOTBUGS_VERSION} spotbugs \
      && wget -nc -O findsecbugs-plugin-1.9.0.jar https://repo1.maven.org/maven2/com/h3xstream/findsecbugs/findsecbugs-plugin/1.9.0/findsecbugs-plugin-1.9.0.jar \
      && echo $SPOTBUGS_VERSION > /opt/spotbugs/version
