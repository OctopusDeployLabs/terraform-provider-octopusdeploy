# terraform-provider-octopusdeploy

A Terraform provider for [Octopus Deploy](https://octopus.com).

It is based on the [go-octopusdeploy](https://github.com/OctopusDeploy/go-octopusdeploy) Octopus Deploy client SDK.

> :warning: This provider is in heavy development. There may be breaking changes.

## Downloading & Installing

As this provider is still under development, you will need to manually download it.

There are compiled binaries for most platforms in [Releases](https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/releases).

To use it, extract the binary for your platform into the same folder as your `.tf` file(s) will be located, then run `terraform init`.

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
