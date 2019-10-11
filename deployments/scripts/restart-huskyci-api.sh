#!/bin/bash
#
# Copyright 2019 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will restart the huskyCI_API container
#

docker stop huskyCI_API > /dev/null
docker rm huskyCI_API > /dev/null
docker-compose -f deployments/docker-compose.yml up -d --build --no-deps api
while true; do
    if [ "$(curl -s -k -L localhost:8888/healthcheck)" = "WORKING" ]; then
        echo "huskyCI_API is UP!"
        break
    else
        echo "Waiting healthcheck..."
    fi
    sleep 15
done
