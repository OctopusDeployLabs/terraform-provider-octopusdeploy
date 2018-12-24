#!/bin/bash
set -eo pipefail

# Download Modules
go mod download

# Lint
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.12.4

./bin/golangci-lint run

# Test
go test -v -timeout 240s ./octopusdeploy/...
