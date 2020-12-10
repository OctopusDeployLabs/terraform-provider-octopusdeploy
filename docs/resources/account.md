---
page_title: "octopusdeploy_account Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages accounts in Octopus Deploy.
---

# Resource `octopusdeploy_account`

This resource manages accounts in Octopus Deploy.



## Schema

### Required

- **name** (String, Required) The name of this resource.

### Optional

- **access_key** (String, Optional) The access key associated with this resource.
- **account_type** (String, Optional) Specifies the type of the account. Valid account types are `AmazonWebServicesAccount`, `AmazonWebServicesRoleAccount`, `AzureServicePrincipal`, `AzureSubscription`, `None`, `SshKeyPair`, `Token`, or `UsernamePassword`.
- **active_directory_endpoint_base_uri** (String, Optional)
- **application_id** (String, Optional) The application ID of this resource.
- **authentication_endpoint** (String, Optional) The authentication endpoint URI for this resource.
- **azure_environment** (String, Optional) The Azure environment associated with this resource. Valid Azure environments are `AzureCloud`, `AzureChinaCloud`, `AzureGermanCloud`, or `AzureUSGovernment`.
- **certificate_data** (String, Optional)
- **certificate_thumbprint** (String, Optional)
- **client_id** (String, Optional)
- **client_secret** (String, Optional)
- **description** (String, Optional) The description of this resource.
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **id** (String, Optional) The unique ID for this resource.
- **password** (String, Optional) The password associated with this resource.
- **private_key_file** (String, Optional)
- **private_key_passphrase** (String, Optional)
- **resource_manager_endpoint** (String, Optional) The resource manager endpoint URI for this resource.
- **secret_key** (String, Optional) The secret key associated with this resource.
- **service_management_endpoint_base_uri** (String, Optional)
- **service_management_endpoint_suffix** (String, Optional)
- **space_id** (String, Optional) The space ID associated with this resource.
- **subscription_id** (String, Optional) The subscription ID of this resource.
- **tenant_id** (String, Optional) The tenant ID of this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.
- **token** (String, Optional) The token of this resource.
- **username** (String, Optional) The username associated with this resource.


