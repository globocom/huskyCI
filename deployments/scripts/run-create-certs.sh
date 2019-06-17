#!/bin/bash
#
# This script will create certs to build the enviroment via Makefile.
#

rm -rf deployments/certs/*
./deployments/scripts/create-certs.sh -m ca -pw "huskyCICert" -t deployments/certs -e 900
./deployments/scripts/create-certs.sh -m server -h dockerapi -pw "huskyCICert" -t deployments/certs -e 365
./deployments/scripts/create-certs.sh -m client -h huskyapi -pw "huskyCICert" -t deployments/certs -e 365
./deployments/scripts/create-certs.sh -m tls -h dockerapi -pw "huskyCICert" -t api -e 365
