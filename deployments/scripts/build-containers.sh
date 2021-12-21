#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will build every securityTest container based on all dockerfiles from huskyCI repository
#

docker build deployments/dockerfiles/bandit/ -t huskyci/bandit:latest
docker build deployments/dockerfiles/brakeman/ -t huskyci/brakeman:latest
docker build deployments/dockerfiles/enry/ -t huskyci/enry:latest
docker build deployments/dockerfiles/gitauthors/ -t huskyci/gitauthors:latest
docker build deployments/dockerfiles/gosec/ -t huskyci/gosec:latest
docker build deployments/dockerfiles/npmaudit/ -t huskyci/npmaudit:latest
docker build deployments/dockerfiles/npmaudit/ -t huskyci/yarnaudit:latest
docker build deployments/dockerfiles/safety/ -t huskyci/safety:latest
docker build deployments/dockerfiles/gitleaks/ -t huskyci/gitleaks:latest
docker build deployments/dockerfiles/spotbugs/ -t huskyci/spotbugs:latest
docker build deployments/dockerfiles/tfsec/ -t huskyci/tfsec:latest
docker build deployments/dockerfiles/securitycodescan/ -t huskyci/securitycodescan:latest