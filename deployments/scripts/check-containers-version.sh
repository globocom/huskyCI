#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will check the version of all securityTests
#


banditVersion=$(docker run --rm huskyci/bandit:latest bandit --version | grep bandit | awk -F " " '{print $2}')
brakemanVersion=$(docker run --rm huskyci/brakeman:latest brakeman --version | awk -F " " '{print $2}')
enryVersion=$(docker run --rm huskyci/enry:latest enry --version)
gitAuthorsVersion=$(docker run --rm huskyci/gitauthors:latest git --version | awk -F " " '{print $3}')
gosecVersion=$(docker run --rm huskyci/gosec:latest gosec --version | grep Version | awk -F " " '{print $2}')
npmAuditVersion=$(docker run --rm huskyci/npmaudit:latest npm audit --version)
yarnAuditVersion=$(docker run --rm huskyci/yarnaudit:latest yarn audit --version )
safetyVersion=$(docker run --rm huskyci/safety:latest safety --version | awk -F " " '{print $3}')
gitleaksVersion=$(docker run --rm huskyci/gitleaks:latest gitleaks --version)
spotbugsVersion=$(docker run --rm huskyci/spotbugs:latest cat /opt/spotbugs/version)
tfsecVersion=$(docker run --rm huskyci/tfsec:latest ./tfsec -v)
securitycodescanVersion=$(docker run --rm huskyci/securitycodescan:latest security-scan | grep tool | awk -F " " '{print $6}')

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
echo "securitycodescanVersion: $securitycodescanVersion"