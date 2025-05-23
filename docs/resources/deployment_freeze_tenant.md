---
page_title: "octopusdeploy_deployment_freeze_tenant Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  
---

# octopusdeploy_deployment_freeze_tenant (Resource)



-> Supported by Octopus Server starting from version 2025.1

## Example Usage

```terraform
# Deployment freeze
resource "octopusdeploy_deployment_freeze" "example" {
  name = "Summer break"
  start = "2024-06-25T00:00:00+10:00"
  end = "2024-06-27T00:00:00+08:00"
}

# Freeze with ids sourced from resources and data sources. 
# Tenants can be sourced from different spaces, a single scope can only reference resources from the same space.

resource "octopusdeploy_deployment_freeze_tenant" "production_freeze" {
  deploymentfreeze_id = octopusdeploy_deployment_freeze.example.id
  tenant_id           = resource.octopusdeploy_tenant.example.id
  project_id          = resource.octopusdeploy_project.example.id
  environment_id      = data.environments.production.id
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `deploymentfreeze_id` (String) The deployment freeze ID associated with this freeze scope.
- `environment_id` (String) The environment ID associated with this freeze scope.
- `project_id` (String) The project ID associated with this freeze scope.
- `tenant_id` (String) The tenant ID associated with this freeze scope.

### Read-Only

- `id` (String) The unique ID for this resource.


