---
layout: "octopusdeploy"
page_title: "Octopus Deploy: machine (Deployment Target)"
---

# Resource: machine (Deployment Target)

Octopus Deploy refers to Machines as [Deployment Targets](https://octopus.com/docs/infrastructure), however the API (and thus Terraform) refers to them as Machines.

## Example Usage

Basic Usage

```hcl
data "octopusdeploy_environment" "staging" {
  name = "Staging"
}

data "octopusdeploy_machinepolicy" "default" {
  name = "Default Machine Policy"
}

resource "octopusdeploy_machine" "testmachine" {
  name                            = "finance-web-01"
  environments                    = ["${data.octopusdeploy_environment.staging.id}"]
  isdisabled                      = false
  machinepolicy                   = "${data.octopusdeploy_machinepolicy.default.id}"
  roles                           = ["Staging"]
  tenanteddeploymentparticipation = "Untenanted"

  endpoint {
    communicationstyle = "TentaclePassive"
    thumbprint         = "81D0FF8B76FC"
    uri                = "https://finance-web-01:10933"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the machine
* `endpoint` - (Required) The configuration of the machine endpoint
    * `communicationstyle` - (Required) Must be one of `None`, `TentaclePassive`, `TentacleActive`, `Ssh`, `OfflineDrop`, `AzureWebApp`, `Ftp`, `AzureCloudService`
    * `proxyid` - (Optional) ID of a defined proxy to use for communication with this machine
    * `thumbprint` - (Required) Thumbprint of the certificate this machine uses (if `communicationstyle` is `None` this should be blank)
    * `uri` - (Required) URI to access this machine (if `endpoint` is `None` this should be blank)
* `environments` - (Required) List of environment IDs to be assigned to this machine
* `isdisabled` - (Required) Whether or not this machine is disabled
* `machinepolicy` - (Required) The ID of the machine policy to be assigned to this machine
* `roles` - (Required) List of the roles to be assigned to this machine
* `tenanteddeploymentparticipation` - (Required) Must be one of `Untenanted`, `TenantedOrUntenanted`, `Tenanted`
* `tenantids` - (Optional) If tenanted, a list of the tenant IDs for this machine
* `tenanttags` - (Optional) If tenanted, a list of the tenant tags for this machine

## Attribute Reference

* `environments` - List of environment IDs this machine is assigned to
* `haslatestcalamari` - Whether or not this machine has the latest Calamari version
* `isdisabled` - Whether or not this machine is disabled
* `isinprocess` - Whether or not this machine is being processed
* `machinepolicy` - The ID of the machine policy assigned to this machine
* `roles` - A list of the roles assigned to this machine
* `status` - The machine status code
* `statussummary` - Plain text description of the machine status
* `tenanteddeploymentparticipation` - One of `Untenanted`, `TenantedOrUntenanted`, `Tenanted`
* `tenantids` - If tenanted, a list of the tenant IDs for this machine
* `tenanttags` -  If tenanted, a list of the tenant tags for this machine