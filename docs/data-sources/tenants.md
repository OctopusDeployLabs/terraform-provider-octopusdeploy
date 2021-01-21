---
page_title: "octopusdeploy_tenants Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing tenants.
---

# Data Source `octopusdeploy_tenants`

Provides information about existing tenants.



## Schema

### Optional

- **cloned_from_tenant_id** (String, Optional) A filter to search for a cloned tenant by its ID.
- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **is_clone** (Boolean, Optional) A filter to search for cloned resources.
- **name** (String, Optional) A filter to search by name.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **project_id** (String, Optional) A filter to search by a project ID.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **tags** (List of String, Optional) A filter to search by a list of tags.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.
- **tenants** (Block List) A list of tenants that match the filter(s). (see [below for nested schema](#nestedblock--tenants))

<a id="nestedblock--tenants"></a>
### Nested Schema for `tenants`

Read-only:

- **cloned_from_tenant_id** (String, Read-only) The ID of the tenant from which this tenant was cloned.
- **description** (String, Read-only) The description of this resource.
- **id** (String, Read-only) The unique ID for this resource.
- **name** (String, Read-only) The name of this resource.
- **project_environment** (Set of Object, Read-only) (see [below for nested schema](#nestedatt--tenants--project_environment))
- **space_id** (String, Read-only) The space ID associated with this resource.
- **tenant_tags** (List of String, Read-only) A list of tenant tags associated with this resource.

<a id="nestedatt--tenants--project_environment"></a>
### Nested Schema for `tenants.project_environment`

- **environments** (List of String)
- **project_id** (String)


