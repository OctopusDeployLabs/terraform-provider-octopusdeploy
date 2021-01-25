---
layout: "octopusdeploy"
page_title: "Octopus Deploy: machine (Deployment Target)"
---

# Data Source: machine (Deployment Target)

Octopus Deploy refers to Machines as [Deployment Targets](https://octopus.com/docs/infrastructure), however the API (and thus Terraform) refers to them as Machines.

## Example Usage

```hcl
resource "octopusdeploy_machine" "testmachine" {
  name = "finance-web-01"
}
```

## Argument Reference

* `name` - (Required) The name of the machine

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
* `endpoint_communicationstyle` - One of `None`, `TentaclePassive`, `TentacleActive`, `Ssh`, `OfflineDrop`, `AzureWebApp`, `Ftp`, `AzureCloudService`
* `endpoint_proxyid` - ID of a defined proxy to use for communication with this machine
* `endpoint_tentacleversiondetails_upgradelocked` - Whether or not this machine tentacle is upgrade locked
* `endpoint_tentacleversiondetails_upgraderequired` - Whether or not this machine tentacle required an upgrade
* `endpoint_tentacleversiondetails_upgradesuggested` - Whether or not this machine tentacle has a suggested ugrade
* `endpoint_tentacleversiondetails_version` - Version number of this machine tentacle
* `endpoint_thumbprint` - Thumbprint of the certificate this machine uses
* `endpoint_uri` - URI to access this machine