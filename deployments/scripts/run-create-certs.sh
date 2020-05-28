#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will create certs to build the enviroment via Makefile.
#
printf 'Generating Certs            ...'
rm -rf .env

if [ ! -f deployments/certs/ca.pem ]; then
    rm -rf deployments/certs/*
    ./deployments/scripts/create-certs.sh -m ca -pw "huskyCICertPassphrase" -t deployments/certs -e 900
    ./deployments/scripts/create-certs.sh -m server -h dockerapi -pw "huskyCICertPassphrase" -t deployments/certs -e 365
    ./deployments/scripts/create-certs.sh -m client -h huskyapi -pw "huskyCICertPassphrase" -t deployments/certs -e 365
fi

echo " done"
