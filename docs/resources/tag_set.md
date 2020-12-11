---
page_title: "octopusdeploy_tag_set Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages tag sets in Octopus Deploy.
---

# Resource `octopusdeploy_tag_set`

This resource manages tag sets in Octopus Deploy.



## Schema

### Required

- **name** (String, Required) The name of this resource.

### Optional

- **description** (String, Optional) The description of this resource.
- **id** (String, Optional) The unique ID for this resource.
- **sort_order** (Number, Optional) The sort order associated with this resource.
- **space_id** (String, Optional) The space ID associated with this resource.
- **tag** (Block List) A list of tags. (see [below for nested schema](#nestedblock--tag))

<a id="nestedblock--tag"></a>
### Nested Schema for `tag`

Required:

- **color** (String, Required)
- **name** (String, Required) The name of this resource.

Optional:

- **canonical_tag_name** (String, Optional)
- **description** (String, Optional) The description of this resource.
- **id** (String, Optional) The unique ID for this resource.
- **sort_order** (Number, Optional)


