---
layout: "octopusdeploy"
page_title: "Provider: Octopus Deploy"
description: |-
  The Octopus Deploy provider provides utilities for configuring and interacting with an Octopus Deploy server.
---

# Octopus Deploy Provider

The Octopus Deploy provider is used to configure resources on an Octopus Deploy server. The provider must be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources and data sources.

## Configure the Provider

### Default Space

Octopus Deploy supports the concept of a Default Space. This is the first space that is automatically created on server setup. If you do not specify a Space when configuring the Octopus Deploy Terraform provider it will use the Default Space.

```hcl
# main.tf

provider "octopusdeploy" {
  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
}
```

### Scoped to a single Space

You can specify a space for the Octopus Deploy Terraform provider to use. If specified, all resources managed by the provider will be scoped to this space. To scope the provider to a space,
simply provide the _name_ of the space (not the space ID).

**Note:** System level resources such as Teams are not support on a Space-scoped provider.

```hcl
# main.tf

provider "octopusdeploy" {
  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
  space   = "Support" # The name of the space
}
```

### Multiple spaces

To manage resources in multiple spaces you can use multiple instances of the provider with [aliases](https://www.terraform.io/docs/configuration/providers.html#alias-multiple-provider-instances) like so:

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
  space   = "Support" # The name of the space
}

provider "octopusdeploy" {
  alias   = "space_product1"

  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
  space   = "Product1" # The name of another space
}

/*
This resource will use the default provider and the default space
*/
resource "octopusdeploy_environment" "Env1" {
  name = "TestEnv1"
}

/*
This resource will use the provider aliased as "space_support"
which is scoped to the space named "Support"
*/
resource "octopusdeploy_environment" "Env2" {
  provider = "octopusdeploy.space_support"
  name     = "TestEnv2"
}

/*
This resource will use the provider aliased as "space_product1"
which is scoped to the space named "Product1"
*/
resource "octopusdeploy_environment" "Env3" {
  provider = "octopusdeploy.space_product1"
  name     = "TestEnv3"
}
```