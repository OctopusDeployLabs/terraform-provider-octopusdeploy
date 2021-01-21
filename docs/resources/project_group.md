---
page_title: "octopusdeploy_project_group Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages project groups in Octopus Deploy.
---

# Resource `octopusdeploy_project_group`

This resource manages project groups in Octopus Deploy.

## Example Usage

```terraform
resource "octopusdeploy_project_group" "example" {
  description  = "The development project group."
  environments = ["Environments-123", "Environments-321"]
  name         = "Development Project Group (OK to Delete)"
}
```

## Schema

### Required

- **name** (String, Required) The name of this resource.

### Optional

- **description** (String, Optional) The description of this resource.
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **id** (String, Optional) The unique ID for this resource.
- **retention_policy_id** (String, Optional)


