#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will tag the version of all securityTests containers
#

rm tmp.txt 2>/dev/null

docker run --rm huskyci/bandit:latest bandit --version > tmp.txt
banditVersion=$(grep bandit tmp.txt | awk -F " " '{print $2}')
docker tag "huskyci/bandit:latest" "huskyci/bandit:$banditVersion"

docker run --rm huskyci/brakeman:latest brakeman --version > tmp.txt
brakemanVersion=$(awk -F " " '{print $2}' < tmp.txt)
docker tag "huskyci/brakeman:latest" "huskyci/brakeman:$brakemanVersion"

docker run --rm huskyci/enry:latest enry --version > tmp.txt
enryVersion=$(cat tmp.txt)
docker tag "huskyci/enry:latest" "huskyci/enry:$enryVersion"

docker run --rm huskyci/gitauthors:latest git --version > tmp.txt
gitAuthorsVersion=$(awk -F " " '{print $3}' < tmp.txt)
docker tag "huskyci/gitauthors:latest" "huskyci/gitauthors:$gitAuthorsVersion"

docker run --rm huskyci/gosec:latest gosec --version > tmp.txt
gosecVersion=$(grep Version tmp.txt | awk -F " " '{print $2}')
docker tag "huskyci/gosec:latest" "huskyci/gosec:$gosecVersion"

docker run --rm huskyci/npmaudit:latest npm audit --version > tmp.txt 
npmAuditVersion=$(cat tmp.txt)
docker tag "huskyci/npmaudit:latest" "huskyci/npmaudit:$npmAuditVersion"

docker run --rm huskyci/yarnaudit:latest yarn audit --version > tmp.txt
yarnAuditVersion=$(cat tmp.txt)
docker tag "huskyci/yarnaudit:latest" "huskyci/yarnaudit:$yarnAuditVersion"

docker run --rm huskyci/safety:latest safety --version > tmp.txt
safetyVersion=$(awk -F " " '{print $3}' < tmp.txt)
docker tag "huskyci/safety:latest" "huskyci/safety:$safetyVersion"

rm tmp.txt 2>/dev/null