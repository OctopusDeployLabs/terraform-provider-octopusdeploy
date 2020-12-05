---
page_title: "octopusdeploy_spaces Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing spaces.
---

# Data Source `octopusdeploy_spaces`

Provides information about existing spaces.

## Example Usage

```terraform
data "octopusdeploy_spaces" "spaces" {
  ids          = ["Spaces-123", "Spaces-321"]
  name         = "Default"
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
```

## Schema

### Optional

- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **name** (String, Optional) A filter to search by name.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.
- **space** (Block List) A list of spaces that match the filter(s). (see [below for nested schema](#nestedblock--space))

<a id="nestedblock--space"></a>
### Nested Schema for `space`

Read-only:

- **description** (String, Read-only) The description of this resource.
- **id** (String, Read-only) The unique identifier for this resource.
- **is_default** (Boolean, Read-only) Specifies if this space is the default space in Octopus.
- **is_task_queue_stopped** (Boolean, Read-only) Specifies the status of the task queue for this space.
- **name** (String, Read-only) The name of this resource.
- **space_managers_team_members** (List of String, Read-only) A list of user IDs designated to be managers of this space.
- **space_managers_teams** (List of String, Read-only) A list of team IDs designated to be managers of this space.


