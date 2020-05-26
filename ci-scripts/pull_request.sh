#!/bin/bash
set -eo pipefail

# Download modules first otherwise linting has errors
go mod download

# Lint

## When in CI, install linter and run from a special path
if [ -n "${CI}" ]; then
    curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.27.0

    ./bin/golangci-lint run
else
    golangci-lint run
fi

# Test

## Enable TF Acceptance Tests
export TF_ACC=1

go test -v -timeout 240s ./octopusdeploy/...
