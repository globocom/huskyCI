#!/bin/bash
#
# This script will generate environment passwords to build huskyCI environment.
#

# huskyCI client default environment variables
HUSKYCI_CLIENT_REPO_URL="https://github.com/globocom/huskyCI.git"
HUSKYCI_CLIENT_REPO_BRANCH="master"
HUSKYCI_CLIENT_API_ADDR="http://localhost:8888"
HUSKYCI_CLIENT_API_USE_HTTPS="false"

# Generating "random" password for certificates
CERT_PASSPHRASE_TMP="certPass$RANDOM$RANDOM"

# Writing password into .env to be used by run_create_cert.sh
echo "export CERT_PASSPHRASE=\"$CERT_PASSPHRASE_TMP\"" > .env

# Adding default envs vars to run be used by make run-client
echo "export HUSKYCI_CLIENT_REPO_URL=\"$HUSKYCI_CLIENT_REPO_URL\"" >> .env
echo "export HUSKYCI_CLIENT_REPO_BRANCH=\"$HUSKYCI_CLIENT_REPO_BRANCH\"" >> .env
echo "export HUSKYCI_CLIENT_API_ADDR=\"$HUSKYCI_CLIENT_API_ADDR\"" >> .env
echo "export HUSKYCI_CLIENT_API_USE_HTTPS=\"$HUSKYCI_CLIENT_API_USE_HTTPS\"" >> .env