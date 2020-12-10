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

- **id** (String, Optional) The unique ID for this resource.
- **tags** (Block List) (see [below for nested schema](#nestedblock--tags))

<a id="nestedblock--tags"></a>
### Nested Schema for `tags`

Required:

- **color** (String, Required)
- **name** (String, Required) The name of this resource.

Optional:

- **canonical_tag_name** (String, Optional)
- **description** (String, Optional) The description of this resource.
- **id** (String, Optional) The unique ID for this resource.
- **sort_order** (Number, Optional)


