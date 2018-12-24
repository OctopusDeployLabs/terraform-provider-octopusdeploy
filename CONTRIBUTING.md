# Contributing

## Go Dependencies

Dependencies are managed using [Go 1.11 Modules](https://github.com/golang/go/wiki/Modules)

## Local Integration Tests

To make development easier, run a local Octopus Deploy server on your machine. You can `vagrant up` [this image](https://github.com/MattHodge/VagrantBoxes/tree/master/OctopusDeployServer) to get a fully working Octopus Deploy Server.

When it comes up, login on [http://localhost:8081](http://localhost:8081) with username `Administrator` and password `OctoVagrant!`.

To get an API to use for local development, go to **Administrator | Profile | My API Keys** and click **New API Key**.

Set the two following environment variables:

```bash
export OCTOPUS_URL=http://localhost:8081/
export OCTOPUS_APIKEY=API-YOUR-API-KEY
```

You can now run integration tests.

## Running Pull Requests Locally

You can locally test pull requests just as the build server would.

- Install [golangci-lint](https://github.com/golangci/golangci-lint)
- Run `./ci-scripts/pull_request.sh`
