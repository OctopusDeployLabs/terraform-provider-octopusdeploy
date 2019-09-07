---
layout: "octopusdeploy"
page_title: "Octopus Deploy: project"
---

# Data Source: octopusdeploy_project

Use this data source to retrieve information about an Octopus Deploy project.

## Example Usage

```hcl
data "octopusdeploy_project" "my_project" {
  name = "My Project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project.

## Attributes Reference

* `description` - A description of the project.

* `lifecycle_id` - The life cycle identifier for the project.

* `project_group_id` - The project group identifer for the project.

* `default_failure_mode` - The failure mode for the project.

* `skip_machine_behavior` - The skip machine behavior for the project.
