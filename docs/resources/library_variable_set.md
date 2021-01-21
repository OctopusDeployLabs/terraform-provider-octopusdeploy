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
- **space_id** (String, Optional) The space ID associated with this resource.
- **template** (Block List) (see [below for nested schema](#nestedblock--template))

### Read-only

- **variable_set_id** (String, Read-only)

<a id="nestedblock--template"></a>
### Nested Schema for `template`

Required:

- **name** (String, Required) The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods. Example: `ServerName`.

Optional:

- **default_value** (String, Optional) A default value for the parameter, if applicable. This can be a hard-coded value or a variable reference.
- **display_settings** (Map of String, Optional) The display settings for the parameter.
- **help_text** (String, Optional) The help presented alongside the parameter input.
- **id** (String, Optional) The unique ID for this resource.
- **label** (String, Optional) The label shown beside the parameter when presented in the deployment process. Example: `Server name`.


