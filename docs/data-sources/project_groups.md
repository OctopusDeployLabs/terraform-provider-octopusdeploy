---
page_title: "octopusdeploy_project_groups Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing project groups.
---

# Data Source `octopusdeploy_project_groups`

Provides information about existing project groups.

## Example Usage

```terraform
data "octopusdeploy_project_groups" "example" {
  ids          = ["ProjectGroups-123", "ProjectGroups-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
```

## Schema

### Optional

- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.
- **project_groups** (Block List) A list of project groups that match the filter(s). (see [below for nested schema](#nestedblock--project_groups))

<a id="nestedblock--project_groups"></a>
### Nested Schema for `project_groups`

Read-only:

- **description** (String, Read-only) The description of this resource.
- **environments** (List of String, Read-only) A list of environment IDs associated with this resource.
- **id** (String, Read-only) The unique ID for this resource.
- **name** (String, Read-only) The name of this resource.
- **retention_policy_id** (String, Read-only)


