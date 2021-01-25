---
layout: "octopusdeploy"
page_title: "Octopus Deploy: lifecycle"
---

# Data Source: octopusdeploy_lifecycle

Use this data source to retrieve information about an Octopus Deploy [lifecycle](https://octopus.com/docs/deployment-process/lifecycles).

## Example Usage

```hcl
data "octopusdeploy_lifecycle" "testing" {
  name = "Testing"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the lifecycle.

## Attributes Reference

* `id` - ID of the lifecycle.

* `description` - A description of the lifecycle.
