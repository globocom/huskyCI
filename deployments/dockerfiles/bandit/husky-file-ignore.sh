#!/bin/sh
#
# Copyright 2020 Globo.com authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# This script will exclude folders and files written in .huskyci from the root of the cloned repository.
#

huskyCIFile=".huskyci"
codePath="/code/"

isHuskyIgnore(){
    line=$1

    if echo "$line" | grep -q "huskyCI-Ignore"; then
        return 0
    fi

    return 1
}

isCommented(){
    line=$1
    commentRegexp='^[[:space:]]*#'

    if echo "$line" | grep -Eq "$commentRegexp"; then
        return 0
    fi

    return 1
}

isEmpty(){
    line=$1

    if [ ! "$line" ]; then
        return 0
    fi

    return 1
}

leavesCodePath(){
    line=$1

    if echo "$line" | grep -qF "../"; then
        return 0
    fi

    return 1
}

wouldRemoveCurrentWorkdir() {
    line=$1

    if echo "$codePath$line" | grep -qF "//"; then
        return 0
    fi

    return 1
}

if [ -f "$huskyCIFile" ]; then

    while IFS= read -r line; do
	    
        if isHuskyIgnore $line || isCommented $line || isEmpty $line || leavesCodePath $line || wouldRemoveCurrentWorkdir $line; then
            continue
        fi

        rm -rf "$codePath$line"

    done < "$huskyCIFile"
fi
