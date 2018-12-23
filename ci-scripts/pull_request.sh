#!/bin/bash
set -eo pipefail

go test -v -timeout 240s ./octopusdeploy/...
