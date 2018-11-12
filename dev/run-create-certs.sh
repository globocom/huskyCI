#!/bin/bash
rm -rf dockers/development/certs/*
./dockers/development/create-certs.sh -m ca -pw yourSecretPassword -t dockers/development/certs -e 900
./dockers/development/create-certs.sh -m server -h dockerapi -pw yourSecretPassword -t dockers/development/certs -e 365
./dockers/development/create-certs.sh -m client -h huskyapi -pw yourSecretPassword -t dockers/development/certs -e 365