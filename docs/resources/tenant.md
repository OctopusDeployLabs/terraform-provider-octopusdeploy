---
page_title: "octopusdeploy_tenant Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages tenants in Octopus Deploy.
---

# Resource `octopusdeploy_tenant`

This resource manages tenants in Octopus Deploy.



## Schema

### Required

- **name** (String, Required) The name of this resource.

### Optional

- **cloned_from_tenant_id** (String, Optional) The ID of the tenant from which this tenant was cloned.
- **description** (String, Optional) The description of this resource.
- **id** (String, Optional) The unique ID for this resource.
- **project_environment** (Block Set) (see [below for nested schema](#nestedblock--project_environment))
- **space_id** (String, Optional) The space ID associated with this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.

<a id="nestedblock--project_environment"></a>
### Nested Schema for `project_environment`

Required:

- **environments** (List of String, Required) A list of environment IDs associated with this tenant through a project.
- **project_id** (String, Required) The ID of the project associated with this tenant.


