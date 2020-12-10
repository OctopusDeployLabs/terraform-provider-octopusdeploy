---
page_title: "octopusdeploy_user_role Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages user roles in Octopus Deploy.
---

# Resource `octopusdeploy_user_role`

This resource manages user roles in Octopus Deploy.



## Schema

### Required

- **name** (String, Required) The name of this resource.

### Optional

- **can_be_deleted** (Boolean, Optional)
- **description** (String, Optional) The description of this resource.
- **granted_space_permissions** (List of String, Optional)
- **granted_system_permissions** (List of String, Optional)
- **id** (String, Optional) The unique ID for this resource.
- **space_permission_descriptions** (List of String, Optional)
- **supported_restrictions** (List of String, Optional)
- **system_permission_descriptions** (List of String, Optional)


