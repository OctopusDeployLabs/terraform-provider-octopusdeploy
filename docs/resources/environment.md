---
page_title: "octopusdeploy_environment Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages environments in Octopus Deploy.
---

# Resource `octopusdeploy_environment`

This resource manages environments in Octopus Deploy.

## Example Usage

```terraform
resource "octopusdeploy_environment" "example" {
  allow_dynamic_infrastructure = false
  description                  = "An environment for the development team."
  name                         = "Development Environment (OK to Delete)"
  use_guided_failure           = false
}
```

## Schema

### Required

- **name** (String, Required) The name of this resource.

### Optional

- **allow_dynamic_infrastructure** (Boolean, Optional)
- **description** (String, Optional) The description of this resource.
- **id** (String, Optional) The unique ID for this resource.
- **use_guided_failure** (Boolean, Optional)

### Read-only

- **sort_order** (Number, Read-only)


