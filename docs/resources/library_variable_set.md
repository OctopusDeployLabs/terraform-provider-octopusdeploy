---
page_title: "octopusdeploy_library_variable_set Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  
---

# Resource `octopusdeploy_library_variable_set`





## Schema

### Required

- **name** (String, Required) The name of this resource.

### Optional

- **description** (String, Optional) The description of this resource.
- **id** (String, Optional) The ID of this resource.
- **space_id** (String, Optional) The space identifier associated with this resource.
- **template** (Block List) (see [below for nested schema](#nestedblock--template))

### Read-only

- **variable_set_id** (String, Read-only)

<a id="nestedblock--template"></a>
### Nested Schema for `template`

Required:

- **name** (String, Required) The name of this resource.

Optional:

- **default_value** (String, Optional)
- **display_settings** (Map of String, Optional)
- **help_text** (String, Optional)
- **id** (String, Optional) The unique identifier for this resource.
- **label** (String, Optional)


