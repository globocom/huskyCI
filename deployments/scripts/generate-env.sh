#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will generate environment passwords to build huskyCI environment.
#

# huskyCI client default environment variables
HUSKYCI_CLIENT_REPO_URL="https://github.com/globocom/huskyCI.git"
HUSKYCI_CLIENT_REPO_BRANCH="poc-golang-gosec"
HUSKYCI_CLIENT_API_ADDR="http://localhost:8888"
HUSKYCI_CLIENT_API_USE_HTTPS="false"

# Adding default envs vars to run be used by make run-client
echo "export HUSKYCI_CLIENT_REPO_URL=\"$HUSKYCI_CLIENT_REPO_URL\"" > .env
echo "export HUSKYCI_CLIENT_REPO_BRANCH=\"$HUSKYCI_CLIENT_REPO_BRANCH\"" >> .env
echo "export HUSKYCI_CLIENT_API_ADDR=\"$HUSKYCI_CLIENT_API_ADDR\"" >> .env
echo "export HUSKYCI_CLIENT_API_USE_HTTPS=\"$HUSKYCI_CLIENT_API_USE_HTTPS\"" >> .env
