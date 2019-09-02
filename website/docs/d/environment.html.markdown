---
layout: "octopusdeploy"
page_title: "Octopus Deploy: environment"
---

# Data Source: octopusdeploy_environment

Use this data source to retrieve information about an Octopus Deploy [environment](https://octopus.com/docs/infrastructure/environments).

## Example Usage

```hcl
data "octopusdeploy_environment" "testing" {
  name = "Testing"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the environment.

## Attributes Reference

* `id` - ID of the environment.

* `description` - A description of the environment.

* `use_guided_failure` - Whether guided failure mode is enabled or not.
