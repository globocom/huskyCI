#!/bin/bash
#
# This script will create certs to build the enviroment via Makefile.
#

. .env
rm -rf deployments/certs/*
./scripts/create-certs.sh -m ca -pw $CERT_PASSPHRASE -t deployments/certs -e 900
./scripts/create-certs.sh -m server -h dockerapi -pw $CERT_PASSPHRASE -t deployments/certs -e 365
./scripts/create-certs.sh -m client -h huskyapi -pw $CERT_PASSPHRASE -t deployments/certs -e 365
