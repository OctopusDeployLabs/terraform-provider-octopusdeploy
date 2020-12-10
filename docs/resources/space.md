---
page_title: "octopusdeploy_space Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages spaces in Octopus Deploy.
---

# Resource `octopusdeploy_space`

This resource manages spaces in Octopus Deploy.

## Example Usage

```terraform
resource "octopusdeploy_space" "example" {
  description                 = "A space for the development team."
  name                        = "Development Team Space"
  is_default                  = false
  is_task_queue_stopped       = false
  space_managers_team_members = ["Users-123", "Users-321"]
  space_managers_teams        = ["teams-everyone"]
}
```

## Schema

### Required

- **name** (String, Required) The name of this resource.

### Optional

- **description** (String, Optional) The description of this resource.
- **id** (String, Optional) The unique ID for this resource.
- **is_default** (Boolean, Optional) Specifies if this space is the default space in Octopus.
- **is_task_queue_stopped** (Boolean, Optional) Specifies the status of the task queue for this space.
- **space_managers_team_members** (List of String, Optional) A list of user IDs designated to be managers of this space.
- **space_managers_teams** (List of String, Optional) A list of team IDs designated to be managers of this space.


