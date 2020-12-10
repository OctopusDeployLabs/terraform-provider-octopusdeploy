---
page_title: "octopusdeploy_token_account Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  
---

# Resource `octopusdeploy_token_account`





## Schema

### Required

- **name** (String, Required) The name of this resource.
- **token** (String, Required) The token of this resource.

### Optional

- **description** (String, Optional) The description of this resource.
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **id** (String, Optional) The unique ID for this resource.
- **space_id** (String, Optional) The space ID associated with this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.


