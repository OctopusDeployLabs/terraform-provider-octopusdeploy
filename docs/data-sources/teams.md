---
page_title: "octopusdeploy_teams Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing users.
---

# Data Source `octopusdeploy_teams`

Provides information about existing users.



## Schema

### Optional

- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **include_system** (Boolean, Optional) A filter to include system teams.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.
- **spaces** (Block List) A list of spaces that match the filter(s). (see [below for nested schema](#nestedblock--spaces))
- **teams** (Block List) A list of teams that match the filter(s). (see [below for nested schema](#nestedblock--teams))

<a id="nestedblock--spaces"></a>
### Nested Schema for `spaces`

Read-only:

- **can_be_deleted** (Boolean, Read-only)
- **can_be_renamed** (Boolean, Read-only)
- **can_change_members** (Boolean, Read-only)
- **can_change_roles** (Boolean, Read-only)
- **description** (String, Read-only)
- **external_security_groups** (List of Object, Read-only) (see [below for nested schema](#nestedatt--spaces--external_security_groups))
- **id** (String, Read-only) The unique ID for this resource.
- **name** (String, Read-only)
- **space_id** (String, Read-only)
- **users** (List of String, Read-only) A list of user IDs designated to be members of this team.

<a id="nestedatt--spaces--external_security_groups"></a>
### Nested Schema for `spaces.external_security_groups`

- **display_id_and_name** (Boolean)
- **display_name** (String)
- **id** (String)



<a id="nestedblock--teams"></a>
### Nested Schema for `teams`

Read-only:

- **can_be_deleted** (Boolean, Read-only)
- **can_be_renamed** (Boolean, Read-only)
- **can_change_members** (Boolean, Read-only)
- **can_change_roles** (Boolean, Read-only)
- **description** (String, Read-only)
- **external_security_groups** (List of Object, Read-only) (see [below for nested schema](#nestedatt--teams--external_security_groups))
- **id** (String, Read-only) The unique ID for this resource.
- **name** (String, Read-only)
- **space_id** (String, Read-only)
- **users** (List of String, Read-only) A list of user IDs designated to be members of this team.

<a id="nestedatt--teams--external_security_groups"></a>
### Nested Schema for `teams.external_security_groups`

- **display_id_and_name** (Boolean)
- **display_name** (String)
- **id** (String)


