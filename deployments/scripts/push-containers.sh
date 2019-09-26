#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will push all securityTests containers
#

rm tmp.txt 2>/dev/null

docker run --rm huskyci/bandit:latest bandit --version > tmp.txt
banditVersion=$(grep bandit tmp.txt | awk -F " " '{print $2}')
docker push "huskyci/bandit:latest" && docker push "huskyci/bandit:$banditVersion"

docker run --rm huskyci/brakeman:latest brakeman --version > tmp.txt
brakemanVersion=$(awk -F " " '{print $2}' < tmp.txt)
docker push "huskyci/brakeman:latest" && docker push "huskyci/brakeman:$brakemanVersion"

docker run --rm huskyci/enry:latest enry --version > tmp.txt
enryVersion=$(cat tmp.txt)
docker push "huskyci/enry:latest" && docker push "huskyci/enry:$enryVersion"

docker run --rm huskyci/gitauthors:latest git --version > tmp.txt
gitAuthorsVersion=$(awk -F " " '{print $3}' < tmp.txt)
docker push "huskyci/gitauthors:latest" && docker push "huskyci/gitauthors:$gitAuthorsVersion"

docker run --rm huskyci/gosec:latest gosec --version > tmp.txt
gosecVersion=$(grep Version tmp.txt | awk -F " " '{print $2}')
docker push "huskyci/gosec:latest" && docker push "huskyci/gosec:$gosecVersion"

docker run --rm huskyci/npmaudit:latest npm audit --version > tmp.txt 
npmAuditVersion=$(cat tmp.txt)
docker push "huskyci/npmaudit:latest" && docker push "huskyci/npmaudit:$npmAuditVersion"

docker run --rm huskyci/yarnaudit:latest yarn audit --version > tmp.txt
yarnAuditVersion=$(cat tmp.txt)
docker push "huskyci/yarnaudit:latest" && docker push "huskyci/yarnaudit:$yarnAuditVersion"

docker run --rm huskyci/safety:latest safety --version > tmp.txt
safetyVersion=$(awk -F " " '{print $3}' < tmp.txt)
docker push "huskyci/safety:latest" && docker push "huskyci/safety:$safetyVersion"

rm tmp.txt 2>/dev/null