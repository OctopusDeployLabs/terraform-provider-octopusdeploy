---
page_title: "octopusdeploy_user_roles Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing user roles.
---

# Data Source `octopusdeploy_user_roles`

Provides information about existing user roles.

## Example Usage

```terraform
data "octopusdeploy_user_roles" "example" {
  ids          = ["UserRoles-123", "UserRoles-321"]
  partial_name = "Administra"
  skip         = 5
  take         = 100
}
```

## Schema

### Optional

- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.
- **user_role** (Block List) A list of user roles that match the filter(s). (see [below for nested schema](#nestedblock--user_role))

<a id="nestedblock--user_role"></a>
### Nested Schema for `user_role`

Read-only:

- **can_be_deleted** (Boolean, Read-only)
- **description** (String, Read-only) The description of this resource.
- **granted_space_permissions** (List of String, Read-only)
- **granted_system_permissions** (List of String, Read-only)
- **id** (String, Read-only) The unique ID for this resource.
- **name** (String, Read-only) The name of this resource.
- **space_permission_descriptions** (List of String, Read-only)
- **supported_restrictions** (List of String, Read-only)
- **system_permission_descriptions** (List of String, Read-only)


