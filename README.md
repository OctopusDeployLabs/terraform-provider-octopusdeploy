# terraform-provider-octopusdeploy

A Terraform provider for [Octopus Deploy](https://octopus.com).

It is based on the [go-octopusdeploy](https://github.com/OctopusDeploy/go-octopusdeploy) Octopus Deploy client SDK.

> :warning: This provider is in heavy development. There may be breaking changes.

## Downloading & Installing

To use this provider you'll need to compile the appropriate binary for your system, place it in the same folder as your `.tf` file(s), then run `terraform init`.

The simplest way to compile a new binary is by using the official [Go Docker Image](https://hub.docker.com/_/golang).
This will enable you to easily produce a binary for [any platform and architecture supported by Go](https://golang.org/doc/install/source#environment).

For example, a 32-bit Windows binary (.exe) is produced with the following command:

```sh
docker run --rm -v "$PWD":/app -w /app -e GOOS=windows -e GOARCH=386 golang go build -v
```

The resulting binary is saved to your `$PWD` as `terraform-provider-octopusdeploy.exe`.

_Note:_ The above command assumes you're running in Bash. If you're using PowerShell, replace `"$PWD"` with `` `"${PWD}`"``. And if you're running plain ol' cmd, use `"%cd%"`.

## Configure the Provider

```hcl
# main.tf

provider "octopusdeploy" {
  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
}
```

## Data Sources

* [octopusdeploy_environment](docs/provider/data_sources/environment.md)
* [octopusdeploy_lifecycle](docs/provider/data_sources/lifecycle.md)

## Provider Resources

* [octopusdeploy_environment](docs/provider/resources/environment.md)
* [octopusdeploy_lifecycle](docs/provider/resources/lifecycle.md)

## Provider Resources (To Be Moved To /docs)

* All other resource documentation is currently [here](docs/to_move_to_provider.md).
