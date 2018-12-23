#!/bin/bash

# Dot Source this script to set the RELEASE_VERSION environment variable

set -eo pipefail

if [ -z "${CIRCLE_BUILD_NUM}" ]; then
    echo "The environment variable CIRCLE_BUILD_NUM is not set. Setting as 999."
    CIRCLE_BUILD_NUM="999"
fi

RELEASE_VERSION="0.0.2-alpha.${CIRCLE_BUILD_NUM}"
echo "Release version is ${RELEASE_VERSION}"

export RELEASE_VERSION=$RELEASE_VERSION
