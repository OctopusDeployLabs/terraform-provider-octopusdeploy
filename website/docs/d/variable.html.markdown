---
layout: "octopusdeploy"
page_title: "Octopus Deploy: variable"
---

# Data Source: variable

## Example Usage

```hcl
data "octopusdeploy_project" "finance" {
    name = "Finance"
}

data "octopusdeploy_environment" "staging" {
    name = "Staging"
}

data "octopusdeploy_variable" "connection_string" {
  project_id = "${data.octopusdeploy_project.finance.id}"
  name       = "SQLConnectionString"
  scope {
    environments = ["${data.octopusdeploy_environment.staging.id}"]
  }
}
```

## Argument Reference

* `project_id` (Required) ID of the Project to assign the variable against.
* `name` - (Required) Name of the variable
* `scope` - (Optional) The scope to apply to this variable. Contains a list of arrays. All are optional:
    * (Optional) `environments`, `machines`, `actions`, `roles`, `channels`, `tenant_tags`

## Attributes reference

* `id` - ID of the variable
* `name` - Name of the variable
* `type` - Type of the variable
* `value` - Value of the variable
* `description` - Description of the variable