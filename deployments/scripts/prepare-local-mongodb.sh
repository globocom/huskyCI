#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will generate mongoDB init JS file. This is done to remove a .JS file from the project
# and avoid huskyCI running Yarn audit and Npm audit
#

cat << EOF > deployments/mongo-init.js
var db = connect("mongodb://localhost/huskyCIDB");

db.createUser(
    {
        user: "huskyCIUser",
        pwd: "huskyCIPassword",
        roles: [{ role: "userAdminAnyDatabase", db: "admin" }]
    }
);

db.createCollection('dockerAPIAddresses');

db.getCollection('dockerAPIAddresses').insertOne({
    currentHostIndex: 0,
    hostList: ["dockerapi"]
});

EOF