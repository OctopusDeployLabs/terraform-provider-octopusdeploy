---
page_title: "octopusdeploy_azure_service_principal Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages Azure service principal accounts in Octopus Deploy.
---

# Resource `octopusdeploy_azure_service_principal`

This resource manages Azure service principal accounts in Octopus Deploy.

## Example Usage

```terraform
resource "octopusdeploy_account" "azure_service_principal_account" {
  account_type    = "AzureServicePrincipal"
  application_id  = "00000000-0000-0000-0000-000000000000"
  name            = "Azure Service Principal Account (OK to Delete)"
  password        = "###########" # required; get from secure environment/store
  subscription_id = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
}
```

## Schema

### Required

- **application_id** (String, Required) The application ID of this resource.
- **name** (String, Required) The name of this resource.
- **password** (String, Required) The password associated with this resource.
- **subscription_id** (String, Required) The subscription ID of this resource.
- **tenant_id** (String, Required) The tenant ID of this resource.

### Optional

- **authentication_endpoint** (String, Optional) The authentication endpoint URI for this resource.
- **azure_environment** (String, Optional) The Azure environment associated with this resource. Valid Azure environments are `AzureCloud`, `AzureChinaCloud`, `AzureGermanCloud`, or `AzureUSGovernment`.
- **description** (String, Optional) The description of this resource.
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **id** (String, Optional) The unique ID for this resource.
- **resource_manager_endpoint** (String, Optional) The resource manager endpoint URI for this resource.
- **space_id** (String, Optional) The space ID associated with this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.


