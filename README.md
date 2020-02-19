# terraform-provider-octopusdeploy

A Terraform provider for [Octopus Deploy](https://octopus.com).

It is based on the [go-octopusdeploy](https://github.com/OctopusDeploy/go-octopusdeploy) Octopus Deploy client SDK.

> :warning: This provider is in heavy development. There may be breaking changes.

## Downloading & Installing

As this provider is still under development, you will need to manually download it.

There are compiled binaries for most platforms in [Releases](https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/releases).

To use it, extract the binary for your platform into the same folder as your `.tf` file(s) will be located, then run `terraform init`.

## Configure the Provider

### Default Space

```hcl
# main.tf

provider "octopusdeploy" {
  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
}
```

### Scoped to a single Space

Simply provide the _name_ of the space (not the space ID)

**Note:** System level resources such as Teams are not support on a Space-scoped provider.

```hcl
# main.tf

provider "octopusdeploy" {
  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
  space   = "Support" // The name of the space
}
```

### Multiple spaces

To manage resources in multiple spaces you currently must use multiple instances of the provider with [aliases](https://www.terraform.io/docs/configuration/providers.html#alias-multiple-provider-instances) like so:

```hcl
# main.tf

provider "octopusdeploy" {
  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
}

provider "octopusdeploy" {
  alias   = "space_support"

  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
  space   = "Support" // The name of the space
}

provider "octopusdeploy" {
  alias   = "space_product1"

  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
  space   = "Product1" // The name of another space
}

// This resource will use the default provider and the default space
resource "octopusdeploy_environment" "Env1" {
  name = "TestEnv1"
}

// This resource will use the provicder aliased as "space_1" which is scoped to "Space-1"
resource "octopusdeploy_environment" "Env2" {
  provider = "octopusdeploy.space_support"
  name     = "TestEnv2"
}

// This resource will use the provider aliased as "space_33" which is scoped to "Space-33"
resource "octopusdeploy_environment" "Env3" {
  provider = "octopusdeploy.space_product1"
  name     = "TestEnv3"
}
```

## Data Sources

* [octopusdeploy_environment](docs/provider/data_sources/environment.md)
* [octopusdeploy_lifecycle](docs/provider/data_sources/lifecycle.md)

## Provider Resources

* [octopusdeploy_environment](docs/provider/resources/environment.md)
* [octopusdeploy_lifecycle](docs/provider/resources/lifecycle.md)
