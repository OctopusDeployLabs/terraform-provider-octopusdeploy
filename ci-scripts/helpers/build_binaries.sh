#!/bin/bash
set -eo pipefail

if [ -z "${RELEASE_VERSION}" ]; then
    echo "The environment variable RELEASE_VERSION needs to be set. Exiting script."
    exit 1
fi

BUILD_PATH_TEMPLATE="build/terraform-provider-octopusdeploy-{{.OS}}_{{.Arch}}-${RELEASE_VERSION}/{{.Dir}}_v${RELEASE_VERSION}"

go get github.com/mitchellh/gox
gox -osarch="linux/amd64" -osarch="linux/386" -osarch="windows/amd64" -osarch="windows/386" -osarch="darwin/amd64" -osarch="darwin/386" -output="${BUILD_PATH_TEMPLATE}"
