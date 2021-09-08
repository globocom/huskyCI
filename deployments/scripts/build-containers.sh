#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will build every securityTest container based on all dockerfiles from huskyCI repository
#

docker build deployments/dockerfiles/bandit/ -t docker-hub.artifactory.globoi.com/huskyci/bandit:latest
docker build deployments/dockerfiles/brakeman/ -t docker-hub.artifactory.globoi.com/huskyci/brakeman:latest
docker build deployments/dockerfiles/enry/ -t docker-hub.artifactory.globoi.com/huskyci/enry:latest
docker build deployments/dockerfiles/gitauthors/ -t docker-hub.artifactory.globoi.com/huskyci/gitauthors:latest
docker build deployments/dockerfiles/gosec/ -t docker-hub.artifactory.globoi.com/huskyci/gosec:latest
docker build deployments/dockerfiles/npmaudit/ -t docker-hub.artifactory.globoi.com/huskyci/npmaudit:latest
docker build deployments/dockerfiles/npmaudit/ -t docker-hub.artifactory.globoi.com/huskyci/yarnaudit:latest
docker build deployments/dockerfiles/safety/ -t docker-hub.artifactory.globoi.com/huskyci/safety:latest
docker build deployments/dockerfiles/gitleaks/ -t docker-hub.artifactory.globoi.com/huskyci/gitleaks:latest
docker build deployments/dockerfiles/spotbugs/ -t docker-hub.artifactory.globoi.com/huskyci/spotbugs:latest
docker build deployments/dockerfiles/tfsec/ -t docker-hub.artifactory.globoi.com/huskyci/tfsec:latest