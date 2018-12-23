#!/bin/bash
set -eo pipefail

if [ -z "${RELEASE_VERSION}" ]; then
    echo "The environment variable RELEASE_VERSION needs to be set. Exiting script."
    exit 1
fi

go get -u github.com/tcnksm/ghr

ghr -u ${CIRCLE_PROJECT_USERNAME} -r ${CIRCLE_PROJECT_REPONAME} -c ${CIRCLE_SHA1} -delete ${RELEASE_VERSION} ./artifacts/
