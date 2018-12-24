#!/bin/bash
set -eo pipefail

ARTIFACTS_PATH="${PWD}/artifacts"
mkdir -p $ARTIFACTS_PATH

for build in $(ls -d build/*)
do
    pushd $build

    RELEASENAME=$(basename $build)
    ZIPNAME="${RELEASENAME}.zip"

    echo "Preparing zip file ${ZIPNAME}"

    # Adding files to zip
    zip $ZIPNAME *

    mv $ZIPNAME $ARTIFACTS_PATH
    popd
done
