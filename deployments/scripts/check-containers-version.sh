#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will check the version of all securityTests
#


banditVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/bandit:latest bandit --version | grep bandit | awk -F " " '{print $2}')
brakemanVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/brakeman:latest brakeman --version | awk -F " " '{print $2}')
enryVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/enry:latest enry --version)
gitAuthorsVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/gitauthors:latest git --version | awk -F " " '{print $3}')
gosecVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/gosec:latest gosec --version | grep Version | awk -F " " '{print $2}')
npmAuditVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/npmaudit:latest npm audit --version)
yarnAuditVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/yarnaudit:latest yarn audit --version )
safetyVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/safety:latest safety --version | awk -F " " '{print $3}')
gitleaksVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/gitleaks:latest gitleaks --version)
spotbugsVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/spotbugs:latest cat /opt/spotbugs/version)
tfsecVersion=$(docker run --rm docker-hub.artifactory.globoi.com/huskyci/tfsec:latest ./tfsec -v)

echo "bandit: $banditVersion"
echo "brakeman: $brakemanVersion"
echo "enry: $enryVersion"
echo "gitauthors: $gitAuthorsVersion"
echo "gosecVersion: $gosecVersion"
echo "npmauditVersion: $npmAuditVersion"
echo "yarnauditVersion: $yarnAuditVersion"
echo "safetyVersion: $safetyVersion"
echo "gitleaksVersion: $gitleaksVersion"
echo "spotbugsVersion: $spotbugsVersion"
echo "tfsecVersion: $tfsecVersion"