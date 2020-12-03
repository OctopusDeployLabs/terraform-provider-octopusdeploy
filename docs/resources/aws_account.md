---
page_title: "octopusdeploy_aws_account Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages AWS accounts in Octopus Deploy.
---

# Resource `octopusdeploy_aws_account`

This resource manages AWS accounts in Octopus Deploy.



## Schema

### Required

- **access_key** (String, Required) The access key associated with this resource.
- **name** (String, Required) The name of this resource.
- **secret_key** (String, Required) The secret key associated with this resource.

### Optional

- **description** (String, Optional) The description of this resource.
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **id** (String, Optional) The unique identifier for this resource.
- **space_id** (String, Optional) The space identifier associated with this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.


