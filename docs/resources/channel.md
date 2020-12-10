---
page_title: "octopusdeploy_channel Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages channels in Octopus Deploy.
---

# Resource `octopusdeploy_channel`

This resource manages channels in Octopus Deploy.



## Schema

### Required

- **name** (String, Required) The name of this resource.
- **project_id** (String, Required) The project ID associated with this channel.

### Optional

- **description** (String, Optional) The description of this resource.
- **id** (String, Optional) The unique ID for this resource.
- **is_default** (Boolean, Optional) Indicates if this is the default channel for the associated project.
- **lifecycle_id** (String, Optional) The lifecycle ID associated with this channel.
- **rules** (Block List) A list of rules associated with this channel. (see [below for nested schema](#nestedblock--rules))
- **space_id** (String, Optional) The space ID associated with this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.

<a id="nestedblock--rules"></a>
### Nested Schema for `rules`

Optional:

- **actions** (List of String, Optional)
- **id** (String, Optional) The unique ID for this resource.
- **tag** (String, Optional)
- **version_range** (String, Optional)


