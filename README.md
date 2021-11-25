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

provider "octopusdeploy" {
  # configuration options
  address    = "https://octopus.example.com"
  api_key    = "API-XXXXXXXXXXXXX"
  space_id   = "Spaces-321"
  space_name = "Development Team"
}
```

There are configuration parameters available for this provider:

* `address` (required; string) the service endpoint of the Octopus REST API
* `api_key` (required; string) the API key to use with the Octopus REST API
* `space_id` (optional; string) the space ID in Octopus Deploy
* `space_name` (optional; string) the space name in Octopus Deploy

If a space ID or name is not specified, the Terraform Provider for Octopus Deploy will assume the default space.

Run `terraform init` to initialize this provider and enable resource management.
