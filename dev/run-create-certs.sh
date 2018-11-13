#!/bin/bash
rm -rf dev/certs/*
./dev/create-certs.sh -m ca -pw yourSecretPassword -t dev/certs -e 900
./dev/create-certs.sh -m server -h dockerapi -pw yourSecretPassword -t dev/certs -e 365
./dev/create-certs.sh -m client -h huskyapi -pw yourSecretPassword -t dev/certs -e 365
