#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will push all securityTests containers
#

banditVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/bandit:latest bandit --version | grep bandit | awk -F " " '{print $2}')
brakemanVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/brakeman:latest brakeman --version | awk -F " " '{print $2}')
enryVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/enry:latest enry --version)
gitAuthorsVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/gitauthors:latest git --version | awk -F " " '{print $3}')
gosecVersion=$(curl -s https://api.github.com/repos/securego/gosec/releases/latest | grep "tag_name" | awk -F '"' '{print $4}')
npmAuditVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/npmaudit:latest npm audit --version)
yarnAuditVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/yarnaudit:latest yarn audit --version )
safetyVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/safety:latest safety --version | awk -F " " '{print $3}')
gitleaksVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/gitleaks:latest gitleaks --version)
spotbugsVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/spotbugs:latest cat /opt/spotbugs/version)
tfsecVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/tfsec:latest ./tfsec -v)

docker tag "docker-hub.artifactory.globoi.com/huskyci/bandit:latest" "docker-hub.artifactory.globoi.com/huskyci/bandit:$banditVersion"
docker tag "docker-hub.artifactory.globoi.com/huskyci/brakeman:latest" "docker-hub.artifactory.globoi.com/huskyci/brakeman:$brakemanVersion"
docker tag "docker-hub.artifactory.globoi.com/huskyci/enry:latest" "docker-hub.artifactory.globoi.com/huskyci/enry:$enryVersion"
docker tag "docker-hub.artifactory.globoi.com/huskyci/gitauthors:latest" "docker-hub.artifactory.globoi.com/huskyci/gitauthors:$gitAuthorsVersion"
docker tag "docker-hub.artifactory.globoi.com/huskyci/gosec:latest" "docker-hub.artifactory.globoi.com/huskyci/gosec:$gosecVersion"
docker tag "docker-hub.artifactory.globoi.com/huskyci/npmaudit:latest" "docker-hub.artifactory.globoi.com/huskyci/npmaudit:$npmAuditVersion"
docker tag "docker-hub.artifactory.globoi.com/huskyci/yarnaudit:latest" "docker-hub.artifactory.globoi.com/huskyci/yarnaudit:$yarnAuditVersion"
docker tag "docker-hub.artifactory.globoi.com/huskyci/safety:latest" "docker-hub.artifactory.globoi.com/huskyci/safety:$safetyVersion"
docker tag "docker-hub.artifactory.globoi.com/huskyci/gitleaks:latest" "docker-hub.artifactory.globoi.com/huskyci/gitleaks:$gitleaksVersion"
docker tag "docker-hub.artifactory.globoi.com/huskyci/spotbugs:latest" "docker-hub.artifactory.globoi.com/huskyci/spotbugs:$spotbugsVersion"
docker tag "docker-hub.artifactory.globoi.com/huskyci/tfsec:latest" "docker-hub.artifactory.globoi.com/huskyci/tfsec:$tfsecVersion"

docker push "docker-hub.artifactory.globoi.com/huskyci/bandit:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/bandit:$banditVersion"
docker push "docker-hub.artifactory.globoi.com/huskyci/brakeman:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/brakeman:$brakemanVersion"
docker push "docker-hub.artifactory.globoi.com/huskyci/enry:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/enry:$enryVersion"
docker push "docker-hub.artifactory.globoi.com/huskyci/gitauthors:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/gitauthors:$gitAuthorsVersion"
docker push "docker-hub.artifactory.globoi.com/huskyci/gosec:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/gosec:$gosecVersion"
docker push "docker-hub.artifactory.globoi.com/huskyci/npmaudit:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/npmaudit:$npmAuditVersion"
docker push "docker-hub.artifactory.globoi.com/huskyci/yarnaudit:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/yarnaudit:$yarnAuditVersion"
docker push "docker-hub.artifactory.globoi.com/huskyci/safety:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/safety:$safetyVersion"
docker push "docker-hub.artifactory.globoi.com/huskyci/gitleaks:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/gitleaks:$gitleaksVersion"
docker push "docker-hub.artifactory.globoi.com/huskyci/spotbugs:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/spotbugs:$spotbugsVersion"
docker push "docker-hub.artifactory.globoi.com/huskyci/tfsec:latest" && docker push "docker-hub.artifactory.globoi.com/huskyci/tfsec:$tfsecVersion"
