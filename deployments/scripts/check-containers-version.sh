#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will check the version of all securityTests
#

rm tmp.txt 2>/dev/null

docker run --rm huskyci/bandit:latest bandit --version > tmp.txt
banditVersion=$(grep bandit tmp.txt | awk -F " " '{print $2}')

docker run --rm huskyci/brakeman:latest brakeman --version > tmp.txt
brakemanVersion=$(awk -F " " '{print $2}' < tmp.txt)

docker run --rm huskyci/enry:latest enry --version > tmp.txt
enryVersion=$(cat tmp.txt)

docker run --rm huskyci/gitauthors:latest git --version > tmp.txt
gitAuthorsVersion=$(awk -F " " '{print $3}' < tmp.txt)

docker run --rm huskyci/gosec:latest gosec --version > tmp.txt
gosecVersion=$(grep Version tmp.txt | awk -F " " '{print $2}')

docker run --rm huskyci/npmaudit:latest npm audit --version > tmp.txt 
npmAuditVersion=$(cat tmp.txt)

docker run --rm huskyci/yarnaudit:latest yarn audit --version > tmp.txt
yarnAuditVersion=$(cat tmp.txt)

docker run --rm huskyci/safety:latest safety --version > tmp.txt
safetyVersion=$(awk -F " " '{print $3}' < tmp.txt)

echo "bandit: $banditVersion"
echo "brakeman: $brakemanVersion"
echo "enry: $enryVersion"
echo "gitauthors: $gitAuthorsVersion"
echo "gosecVersion: $gosecVersion"
echo "npmauditVersion: $npmAuditVersion"
echo "yarnauditVersion: $yarnAuditVersion"
echo "safetyVersion: $safetyVersion"

rm tmp.txt 2>/dev/null