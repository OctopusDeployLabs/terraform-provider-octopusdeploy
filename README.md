# terraform-provider-octopusdeploy
A Terraform provider for [Octopus Deploy].

[![Build status](https://ci.appveyor.com/api/projects/status/5t5gbqjyl8hpou52?svg=true)](https://ci.appveyor.com/project/MattHodge/go-octopusdeploy)

> :warning: This provider is in heavy development. It is not production ready yet.

# Go Dependencies
* Dependencies are managed using [dep](https://golang.github.io/dep/docs/new-project.html)

```bash
# Vendor new modules
dep ensure
```

# Example

```json
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
