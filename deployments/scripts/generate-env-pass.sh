#!/bin/bash
#
# This script will generate environment passwords to build the environment.
#

# huskyCI client default vars
HUSKYCI_REPO_URL="https://github.com/globocom/huskyci.git"
HUSKYCI_REPO_BRANCH="master"
HUSKYCI_API="https://localhost:8888"

# Generating "random" passwords
CERT_PASSPHRASE_TMP="certPass$RANDOM$RANDOM"
MONGO_DATABASE_USERNAME_TMP="huskyUser$RANDOM$RANDOM"
MONGO_DATABASE_PASSWORD_TMP="huskyPass$RANDOM$RANDOM"

# Writing passwords into dockers.env file to be used by docker compose
echo "MONGO_DATABASE_USERNAME=$MONGO_DATABASE_USERNAME_TMP" > deployments/dockers.env
echo "MONGO_DATABASE_PASSWORD=$MONGO_DATABASE_PASSWORD_TMP" >> deployments/dockers.env

# Writing passwords into .env to be used by run_create_cert.sh and to send to STDOUT
echo "export CERT_PASSPHRASE=\"$CERT_PASSPHRASE_TMP\"" > .env
echo "export MONGO_DATABASE_USERNAME=\"$MONGO_DATABASE_USERNAME_TMP\"" >> .env
echo "export MONGO_DATABASE_PASSWORD=\"$MONGO_DATABASE_PASSWORD_TMP\"" >> .env

# Adding default envs vars to run client
echo "export HUSKYCI_REPO_URL=\"$HUSKYCI_REPO_URL\"" >> .env
echo "export HUSKYCI_REPO_BRANCH=\"$HUSKYCI_REPO_BRANCH\"" >> .env
echo "export HUSKYCI_API=\"$HUSKYCI_API\"" >> .env

# Preparing script to create mongoDB default user
cat << EOF > deployments/mongo-init.js
var db = connect("mongodb://localhost/huskyDB");

db.createUser(
    {
        user: "${MONGO_DATABASE_USERNAME_TMP}",
        pwd: "${MONGO_DATABASE_PASSWORD_TMP}",
        roles: [{ role: "userAdminAnyDatabase", db: "admin" }]
    }
);
EOF
