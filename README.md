# terraform-provider-octopusdeploy
A Terraform provider for [Octopus Deploy](https://octopus.com).

[![Build status](https://ci.appveyor.com/api/projects/status/a5ejcududsoug94e/branch/master?svg=true)](https://ci.appveyor.com/project/MattHodge/terraform-provider-octopusdeploy/branch/master)

> :warning: This provider is in heavy development. It is not production ready yet.

# Go Dependencies
* Dependencies are managed using [dep](https://golang.github.io/dep/docs/new-project.html)

```bash
# Vendor new modules
dep ensure
```

# Provider Resources

Resource Name | Description
--- | ---
`octopusdeploy_project` | Create an Octopus Deploy project
`octopusdeploy_project_group` | Create an Octopus Deploy project group


# Example

```hcl
# main.tf

provider "octopusdeploy" {
  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
}

resource "octopusdeploy_project_group" "test_projectgroup" {
  description = "The Best Team"
  name        = "Team #1"
}

resource "octopusdeploy_project" "test_project" {
  description      = "An example description"
  lifecycle_id     = "Lifecycles-1"
  name             = "My Octopus Deploy Project"
  project_group_id = "${octopusdeploy_project_group.test_projectgroup.id}"
}

```
