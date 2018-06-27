# terraform-provider-octopusdeploy
A Terraform provider for [Octopus Deploy](https://octopus.com).

> :warning: This provider is in heavy development. It is not production ready yet.

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
