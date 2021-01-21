---
page_title: "octopusdeploy_aws_account Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages AWS accounts in Octopus Deploy.
---

# Resource `octopusdeploy_aws_account`

This resource manages AWS accounts in Octopus Deploy.

## Example Usage

```terraform
resource "octopusdeploy_account" "amazon_web_services_account" {
  access_key   = "access-key"
  account_type = "AmazonWebServicesAccount"
  name         = "AWS Account (OK to Delete)"
  secret_key   = "###########" # required; get from secure environment/store
}
```

## Schema

### Required

- **access_key** (String, Required) The access key associated with this AWS account.
- **name** (String, Required) The name of this AWS account.
- **secret_key** (String, Required) The secret key associated with this resource.

### Optional

- **description** (String, Optional) A user-friendly description of this AWS account.
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **id** (String, Optional) The ID of this resource.
- **space_id** (String, Optional) The space ID associated with this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.


