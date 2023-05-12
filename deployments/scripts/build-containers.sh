#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will build every securityTest container based on all dockerfiles from huskyCI repository
#

docker buildx build --platform linux/amd64 deployments/dockerfiles/bandit/ -t huskyci/bandit:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/brakeman/ -t huskyci/brakeman:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/enry/ -t huskyci/enry:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/gitauthors/ -t huskyci/gitauthors:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/gosec/ -t huskyci/gosec:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/npmaudit/ -t huskyci/npmaudit:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/npmaudit/ -t huskyci/yarnaudit:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/safety/ -t huskyci/safety:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/gitleaks/ -t huskyci/gitleaks:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/spotbugs/ -t huskyci/spotbugs:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/tfsec/ -t huskyci/tfsec:latest
docker buildx build --platform linux/amd64 deployments/dockerfiles/securitycodescan/ -t huskyci/securitycodescan:latest
