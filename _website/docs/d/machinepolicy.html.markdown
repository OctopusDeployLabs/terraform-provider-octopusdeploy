---
layout: "octopusdeploy"
page_title: "Octopus Deploy: machinepolicy"
---

# Data Source: machinepolicy

[Machine policies](https://octopus.com/docs/infrastructure/machine-policies) are groups of settings that can be applied to Tentacle and SSH endpoints to modify their behavior.

Currently the Octopus terraform provider only provides machine policies as a data provider, as there are places elsewhere in
the provider that the IDs of machine policies need to be referenced.

## Example Usage

```hcl
data "octopusdeploy_machinepolicy" "default" {
  name = "Default Machine Policy"
}
```

## Argument Reference

* `name` - (Required) The name of the machine policy

## Attributes Reference

* `name` - The name of the machine policy
* `description` - The description of the machine policy
* `isdefault` - Whether or not this machine policy is the default policy