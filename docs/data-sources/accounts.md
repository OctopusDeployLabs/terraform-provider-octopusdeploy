---
page_title: "octopusdeploy_accounts Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing accounts.
---

# Data Source `octopusdeploy_accounts`

Provides information about existing accounts.



## Schema

### Optional

- **account_type** (String, Optional) A filter to search by a list of account types.  Valid account types are `AmazonWebServicesAccount`, `AmazonWebServicesRoleAccount`, `AzureServicePrincipal`, `AzureSubscription`, `None`, `SshKeyPair`, `Token`, or `UsernamePassword`.
- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **accounts** (List of Object, Read-only) A list of accounts that match the filter(s). (see [below for nested schema](#nestedatt--accounts))
- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.

<a id="nestedatt--accounts"></a>
### Nested Schema for `accounts`

- **access_key** (String)
- **account_type** (String)
- **active_directory_endpoint_base_uri** (String)
- **application_id** (String)
- **authentication_endpoint** (String)
- **azure_environment** (String)
- **certificate_data** (String)
- **certificate_thumbprint** (String)
- **client_id** (String)
- **client_secret** (String)
- **description** (String)
- **environments** (List of String)
- **id** (String)
- **name** (String)
- **password** (String)
- **private_key_file** (String)
- **private_key_passphrase** (String)
- **resource_manager_endpoint** (String)
- **secret_key** (String)
- **service_management_endpoint_base_uri** (String)
- **service_management_endpoint_suffix** (String)
- **space_id** (String)
- **subscription_id** (String)
- **tenant_id** (String)
- **tenant_tags** (List of String)
- **tenanted_deployment_participation** (String)
- **tenants** (List of String)
- **token** (String)
- **username** (String)


