#!/bin/bash
set -eo pipefail

# Clean the directories
ci-scripts/helpers/clean.sh

# Load the release version environment variable
. ci-scripts/helpers/get_release_version.sh

# Build the binaries
ci-scripts/helpers/build_binaries.sh

# Create zip files
ci-scripts/helpers/create_artifacts.sh
