---
page_title: "Octopus Deploy Provider Configuration"
subcategory: "Guides"
---

# Octopus Deploy Provider Configuration

## Example usage

### API Key

`main.tf`

```hcl
terraform {
    required_providers {
        octopusdeploy = {
            source = "OctopusDeployLabs/octopusdeploy"
        }
    }
}

provider "octopusdeploy" {
  address       = "https://octopus-cloud-or-server-uri"
  api_key       = "octous_deploy_api_key"
  space_id      = "..."
}
```

### Access Token (via Environment Variable)
OIDC Access Tokens are short-lived and typically generated per-run of an automated pipeline, such as GitHub Actions.
If you use the Access Token approach, we recommend sourcing the token from environment variable.

The environment variable fallback values that the Terraform Provider search for correspond to the values that pipeline steps like our [GitHub Login action](https://github.com/OctopusDeploy/login?tab=readme-ov-file#outputs) set in the pipeline context, so the provider will automatically pick up the value from environment variable.

`main.tf`

```hcl
terraform {
    required_providers {
        octopusdeploy = {
            source = "OctopusDeployLabs/octopusdeploy"
        }
    }
}

provider "octopusdeploy" {
  space_id      = "..."
}
```

## Schema

### Required
* `address` (String) The Octopus Deploy server URL.

and one of either
* `api_key` (String) The Octopus Deploy server API key.

OR
* `access_token` (String) The OIDC Access Token from an OIDC exchange.

### Optional
* `space_id` (String) The ID of the space to create the resources in.

**If `space_id` is not specified the default space will be used.**

### Environment Variable fallback
The following priority order will be used to calculate the final value for these configuration items:

| Configuration Item | Priority Order                                                                                   |
|--------------------|--------------------------------------------------------------------------------------------------|
| `address`          | 1. Provider Configuration Block <br /> 2. env: `OCTOPUS_URL`                                     |
| `api_key`          | 1. Provider Configuration Block <br /> 2. env: `OCTOPUS_APIKEY` <br /> 3. env: `OCTOPUS_API_KEY` |
| `access_token`     | 1. Provider Configuration Block <br /> 2. env: `OCTOPUS_ACCESS_TOKEN`                            |
