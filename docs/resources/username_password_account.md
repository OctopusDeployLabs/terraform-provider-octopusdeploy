---
page_title: "octopusdeploy_username_password_account Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages username-password accounts in Octopus Deploy.
---

# Resource `octopusdeploy_username_password_account`

This resource manages username-password accounts in Octopus Deploy.



## Schema

### Required

- **name** (String, Required) The name of this resource.
- **username** (String, Required) The username associated with this resource.

### Optional

- **description** (String, Optional) The description of this resource.
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **id** (String, Optional) The unique ID for this resource.
- **password** (String, Optional) The password associated with this resource.
- **space_id** (String, Optional) The space ID associated with this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.


