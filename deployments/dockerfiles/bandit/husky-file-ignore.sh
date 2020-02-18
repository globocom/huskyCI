#!/bin/sh
#
# Copyright 2020 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will exclude folders and files written in .huskyci from the root of the cloned repository.
#

huskyCIFile=".huskyci"

if [ -f "$huskyCIFile" ]; then
    while IFS= read -r line; do
	commentRegexp='^[[:space:]]*#'
        if echo "$line" | grep -q "huskyCI-Ignore" || echo "$line" | grep -Eq "$commentRegexp" || [ ! "$line" ]; then
            continue
        fi

        rm -rf $line

    done < "$huskyCIFile"
fi
