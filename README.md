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

# Example

```hcl
# main.tf

provider "octopusdeploy" {
  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
}

resource "octopusdeploy_project" "test_project" {
  name           = "My Octopus Deploy Project"
  lifecycleid    = "Lifecycles-1"
  projectgroupid = "ProjectGroups-1"
  description    = "An example description"
}
```
