---
page_title: "2. Octopus Deploy Provider Configuration"
subcategory: "Guides"
---

# 2. Octopus Deploy Provider Configuration

## Example usage

`main.tf`

```hcl
terraform {
    required_providers {
        octopusdeploy = {
            source = OctopusDeployLabs/octopusdeploy
        }
    }
}

provider "octopusdeploy" {
  address       = "https://octopus-cloud-or-server-uri"
  api_key       = "octous_deploy_api_key"
  space_id      = "..."
}
```

## Schema

### Required
* `address` (String) The Octopus Deploy server URL. This can also be set using the `OCTOPUS_URL` environment variable.
* `api_key` (String) The Octopus Deploy server API key. This can also be set using the `OCTOPUS_APIKEY` environment variable.

### Optional
* `space_id` (String) The ID of the space to create the resources in.

**If `space_id` is not specified the default space will be used.**