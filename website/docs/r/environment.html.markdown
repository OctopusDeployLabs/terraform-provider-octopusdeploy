---
layout: "octopusdeploy"
page_title: "Octopus Deploy: environment"
---

# Resource: octopusdeploy_environment

Use this resource allows the creation of Octopus Deploy [environment](https://octopus.com/docs/infrastructure/environments).

Environments help you organize your deployment targets.

## Example Usage

```hcl
resource "octopusdeploy_environment" "staging" {
    name               = "Staging"
    description        = "Staging environment"
    use_guided_failure = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the environment.

* `description` - (Optional) Description of the environment.

* `use_guided_failure` - (Optional) Use guided failures for this environment. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the environment.
