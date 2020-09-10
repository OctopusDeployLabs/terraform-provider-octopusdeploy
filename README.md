# terraform-provider-octopusdeploy
> :warning: This is a community project under development. Please raise a GitHub issue for any problems or feature requests.

![Run integration tests against Octopus in a Docker container](https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/workflows/Run%20integration%20tests%20against%20Octopus%20in%20a%20Docker%20container/badge.svg)

A Terraform provider for [Octopus Deploy](https://octopus.com).

It is based on the [go-octopusdeploy](https://github.com/OctopusDeploy/go-octopusdeploy) Octopus Deploy client SDK.

## Testing

A GitHub Action Workflow has been added to this project which initializes an instance of Octopus Deploy and runs the tests against it. These same tests can be run in a forked repository.

The GitHub Action Workflow can be found [here](https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/actions?query=workflow%3A%22Run+integration+tests+against+Octopus+in+a+Docker+container%22)

## Downloading & Installing
We are actively working with Hashicorp to join the partner program so you don't need to do the following manual steps. 

As this provider is still under development, you will need to manually download it.

You can find the most recent compiled binary here [Releases](https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/releases/tag/v0.6.0).

To use it the binary:

1. The first command you will want to run is go get to pull down the executable

`go get github.com/OctopusDeploy/terraform-provider-octopusdeploy`

Once the executable is pulled down, it'll automatically go into the ~/go directory on Linux/MacOS or the go directory on the home folder in Windows. Three folders will be shown in the go directory:

* bin
* pkg
* src

Typically the executable is in the *bin directory.

2. cd into ~/go/bin

3. Switching gears for a moment - in the directory where you want the Terraform configuration files to exist to use the Octopus Deploy Terraform provider, create the directory .terraform/plugins/OS plugin. The OS Plugin will be different based on OS, so for example, MacOS would look like .terraform/plugins/darwin_amd64

4. Copy the terraform-provider-octopusdeploy into `.`terraform/plugins/OS plugin`

You should know be able to initialize the environment.

## Configure the Provider

### Provider
The provider will always be used per the code below. You will need to specify the Octopus Deploy server address, the API key (which should be passed in securely), and the space name.

```hcl
# main.tf

provider "octopusdeploy" {
  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
  space   = "Space Name
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
