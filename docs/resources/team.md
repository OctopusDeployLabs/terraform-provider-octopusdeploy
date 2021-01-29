---
page_title: "octopusdeploy_team Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages teams in Octopus Deploy.
---

# Resource `octopusdeploy_team`

This resource manages teams in Octopus Deploy.



## Schema

### Required

- **name** (String, Required)

### Optional

- **can_be_deleted** (Boolean, Optional)
- **can_be_renamed** (Boolean, Optional)
- **can_change_members** (Boolean, Optional)
- **can_change_roles** (Boolean, Optional)
- **description** (String, Optional)
- **external_security_groups** (Block List) (see [below for nested schema](#nestedblock--external_security_groups))
- **id** (String, Optional) The unique ID for this resource.
- **users** (List of String, Optional) A list of user IDs designated to be members of this team.

### Read-only

- **space_id** (String, Read-only)

<a id="nestedblock--external_security_groups"></a>
### Nested Schema for `external_security_groups`

Optional:

- **display_id_and_name** (Boolean, Optional)
- **display_name** (String, Optional)
- **id** (String, Optional) The unique ID for this resource.


