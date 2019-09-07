---
layout: "octopusdeploy"
page_title: "Octopus Deploy: project_group"
---

## Resource: Project Groups

[Project groups](https://octopus.com/docs/deployment-process/projects#project-group) are a way of organizing your projects.

### Example Usage

Basic usage:

```hcl
resource "octopusdeploy_project_group" "finance" {
  description = "Financial Applications"
  name        = "Finance"
}
```

Basic usage with ID export used to create project:

```hcl
resource "octopusdeploy_project_group" "finance" {
  description = "Financial Applications"
  name        = "Finance"
}

resource "octopusdeploy_project" "billing_service" {
  description      = "The Finance teams billing service"
  lifecycle_id     = "Lifecycles-1"
  name             = "Billing Service"
  project_group_id = "${octopusdeploy_project_group.finance.id}"
}
```

### Argument Reference
* `description` - (Optional) Description of the project group
* `name` - (Required) Name of the project group

### Attributes Reference
* `id` - The ID of the project group