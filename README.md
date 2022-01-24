# 🐙 Terraform Provider for Octopus Deploy

This repository contains the source code for the Terraform Provider for [Octopus Deploy](https://octopus.com). It supports provisioning/configuring of Octopus Deploy instances via [Terraform](https://www.terraform.io/). Documentation and guides for using this provider are located on the Terraform Registry: [Documentation](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs).

## 🪄 Installation and Configuration

The Terraform Provider for Octopus Deploy is available via the Terraform Registry: [OctopusDeployLabs/octopusdeploy](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy). To install this provider, copy and paste this code into your Terraform configuration:

```hcl
terraform {
  required_providers {
    octopusdeploy = {
      source = "OctopusDeployLabs/octopusdeploy"
      version = "version-number" # example: 0.7.62
    }
  }
}

data "octopusdeploy_space" "space" {
  provider = octopusdeploy.unscoped
  name     = "Development Team"
}

provider "octopusdeploy" {
  # configuration options
  address    = "https://octopus.example.com"     # (required; string) the service endpoint of the Octopus REST API
  api_key    = "API-XXXXXXXXXXXXX"               # (required; string) the API key to use with the Octopus REST API
  space_id   = data.octopusdeploy_space.space.id # (optional; string) the space ID in Octopus Deploy
}
```

If `space_id` is not specified the Terraform Provider for Octopus Deploy will assume the default space.

### Environment Variables

You can provide your Octopus Server configuration via the `OCTOPUS_URL` and `OCTOPUS_APIKEY` environment variables, representing your Octopus Server address and API Key, respectively.

```hcl
provider "octopusdeploy" {}
```

Run `terraform init` to initialize this provider and enable resource management.

## 🛠 Build Instructions

A build of this Terraform Provider can be created using the [Makefile](https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/blob/master/Makefile) provided in the source:

```
% make build -f Makefile
```

This will generate a binary that will be installed to the local plugins folder. Once installed, the provider may be used through the following configuration:

```hcl
terraform {
  required_providers {
    octopusdeploy = {
      source  = "octopus.com/com/octopusdeploy"
      version = "0.7.64"
    }
  }
}

provider "octopusdeploy" {
  address  = # address
  api_key  = # API key
  space_id = # space ID
}
```
