---
page_title: "octopusdeploy_azure_subscription_account Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages Azure subscription accounts in Octopus Deploy.
---

# Resource `octopusdeploy_azure_subscription_account`

This resource manages Azure subscription accounts in Octopus Deploy.



## Schema

### Required

- **management_endpoint** (String, Required)
- **name** (String, Required) The name of this resource.
- **storage_endpoint_suffix** (String, Required) The storage endpoint suffix associated with this Azure subscription account.
- **subscription_id** (String, Required) The subscription ID of this resource.

### Optional

- **azure_environment** (String, Optional) The Azure environment associated with this resource. Valid Azure environments are `AzureCloud`, `AzureChinaCloud`, `AzureGermanCloud`, or `AzureUSGovernment`.
- **certificate** (String, Optional)
- **certificate_thumbprint** (String, Optional)
- **description** (String, Optional) The description of this resource.
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **id** (String, Optional) The ID of this resource.
- **space_id** (String, Optional) The space ID associated with this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.


