#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will tag the version of all securityTests containers
#

banditVersion=$(docker run --rm huskyci/bandit:latest bandit --version | grep bandit | awk -F " " '{print $2}')
brakemanVersion=$(docker run --rm huskyci/brakeman:latest brakeman --version | awk -F " " '{print $2}')
enryVersion=$(docker run --rm huskyci/enry:latest enry --version)
gitAuthorsVersion=$(docker run --rm huskyci/gitauthors:latest git --version | awk -F " " '{print $3}')
gosecVersion=$(docker run --rm huskyci/gosec:latest gosec --version | grep Version | awk -F " " '{print $2}')
npmAuditVersion=$(docker run --rm huskyci/npmaudit:latest npm audit --version)
yarnAuditVersion=$(docker run --rm huskyci/yarnaudit:latest yarn audit --version )
safetyVersion=$(docker run --rm huskyci/safety:latest safety --version | awk -F " " '{print $3}')

docker tag "huskyci/bandit:latest" "huskyci/bandit:$banditVersion"
docker tag "huskyci/brakeman:latest" "huskyci/brakeman:$brakemanVersion"
docker tag "huskyci/enry:latest" "huskyci/enry:$enryVersion"
docker tag "huskyci/gitauthors:latest" "huskyci/gitauthors:$gitAuthorsVersion"
docker tag "huskyci/gosec:latest" "huskyci/gosec:$gosecVersion"
docker tag "huskyci/npmaudit:latest" "huskyci/npmaudit:$npmAuditVersion"
docker tag "huskyci/yarnaudit:latest" "huskyci/yarnaudit:$yarnAuditVersion"
docker tag "huskyci/safety:latest" "huskyci/safety:$safetyVersion"