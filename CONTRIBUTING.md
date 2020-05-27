# Contributing

## Testing and Linting

The GitHub action in the workflow file `.github/test.yml` creates an Octopus instance and then runs linting and tests
against it. A PR will only be accepted if it doesn't break any existing passing tests, if all linting rules also pass,
and if the new feature has additional tests.

GitHub actions are available in forked repos, but must be manually enabled.

## Go Dependencies

Dependencies are managed using [Go 1.11 Modules](https://github.com/golang/go/wiki/Modules)

## Local Integration Tests

To make development easier, run a local Octopus Deploy server on your machine. The Docker Compose file at
`tests/docker-compose.yml` will create a test environment. Set the `OCTOPUS_VERSION` environment variable to a valid
`octopusdeploy/octopusdeploy` [image tag](https://hub.docker.com/r/octopusdeploy/octopusdeploy).

When it comes up, login on [http://localhost:8080](http://localhost:8080) with username `admin` and password `Password01!`.

To get an API to use for local development, go to **Administrator | Profile | My API Keys** and click **New API Key**.

Set the two following environment variables:

```bash
export OCTOPUS_URL=http://localhost:8080/
export OCTOPUS_APIKEY=API-YOUR-API-KEY
```

You can now run integration tests.

## Running Pull Requests Locally

You can locally test pull requests just as the build server would.

- Install [golangci-lint](https://github.com/golangci/golangci-lint)
- Run `./ci-scripts/pull_request.sh`
