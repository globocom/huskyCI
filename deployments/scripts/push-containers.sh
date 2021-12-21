#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will push all securityTests containers
#

banditVersion=$(docker run --rm huskyci/bandit:latest bandit --version | grep bandit | awk -F " " '{print $2}')
brakemanVersion=$(docker run --rm huskyci/brakeman:latest brakeman --version | awk -F " " '{print $2}')
enryVersion=$(docker run --rm huskyci/enry:latest enry --version)
gitAuthorsVersion=$(docker run --rm huskyci/gitauthors:latest git --version | awk -F " " '{print $3}')
gosecVersion=$(curl -s https://api.github.com/repos/securego/gosec/releases/latest | grep "tag_name" | awk -F '"' '{print $4}')
npmAuditVersion=$(docker run --rm huskyci/npmaudit:latest npm audit --version)
yarnAuditVersion=$(docker run --rm huskyci/yarnaudit:latest yarn audit --version )
safetyVersion=$(docker run --rm huskyci/safety:latest safety --version | awk -F " " '{print $3}')
gitleaksVersion=$(docker run --rm huskyci/gitleaks:latest gitleaks --version)
spotbugsVersion=$(docker run --rm huskyci/spotbugs:latest cat /opt/spotbugs/version)
tfsecVersion=$(docker run --rm huskyci/tfsec:latest ./tfsec -v)
securitycodescanVersion=$(docker run --rm huskyci/securitycodescan:latest security-scan | grep tool | awk -F " " '{print $6}')

docker tag "huskyci/bandit:latest" "huskyci/bandit:$banditVersion"
docker tag "huskyci/brakeman:latest" "huskyci/brakeman:$brakemanVersion"
docker tag "huskyci/enry:latest" "huskyci/enry:$enryVersion"
docker tag "huskyci/gitauthors:latest" "huskyci/gitauthors:$gitAuthorsVersion"
docker tag "huskyci/gosec:latest" "huskyci/gosec:$gosecVersion"
docker tag "huskyci/npmaudit:latest" "huskyci/npmaudit:$npmAuditVersion"
docker tag "huskyci/yarnaudit:latest" "huskyci/yarnaudit:$yarnAuditVersion"
docker tag "huskyci/safety:latest" "huskyci/safety:$safetyVersion"
docker tag "huskyci/gitleaks:latest" "huskyci/gitleaks:$gitleaksVersion"
docker tag "huskyci/spotbugs:latest" "huskyci/spotbugs:$spotbugsVersion"
docker tag "huskyci/tfsec:latest" "huskyci/tfsec:$tfsecVersion"
docker tag "huskyci/securitycodescan:latest" "huskyci/securitycodescan:$securitycodescanVersion"

docker push "huskyci/bandit:latest" && docker push "huskyci/bandit:$banditVersion"
docker push "huskyci/brakeman:latest" && docker push "huskyci/brakeman:$brakemanVersion"
docker push "huskyci/enry:latest" && docker push "huskyci/enry:$enryVersion"
docker push "huskyci/gitauthors:latest" && docker push "huskyci/gitauthors:$gitAuthorsVersion"
docker push "huskyci/gosec:latest" && docker push "huskyci/gosec:$gosecVersion"
docker push "huskyci/npmaudit:latest" && docker push "huskyci/npmaudit:$npmAuditVersion"
docker push "huskyci/yarnaudit:latest" && docker push "huskyci/yarnaudit:$yarnAuditVersion"
docker push "huskyci/safety:latest" && docker push "huskyci/safety:$safetyVersion"
docker push "huskyci/gitleaks:latest" && docker push "huskyci/gitleaks:$gitleaksVersion"
docker push "huskyci/spotbugs:latest" && docker push "huskyci/spotbugs:$spotbugsVersion"
docker push "huskyci/tfsec:latest" && docker push "huskyci/tfsec:$tfsecVersion"
docker push "huskyci/securitycodescan:latest" && docker push "huskyci/securitycodescan:$securitycodescanVersion"
