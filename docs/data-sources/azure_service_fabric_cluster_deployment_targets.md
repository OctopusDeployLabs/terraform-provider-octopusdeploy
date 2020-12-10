---
page_title: "octopusdeploy_azure_service_fabric_cluster_deployment_targets Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing Azure service fabric cluster deployment targets.
---

# Data Source `octopusdeploy_azure_service_fabric_cluster_deployment_targets`

Provides information about existing Azure service fabric cluster deployment targets.



## Schema

### Optional

- **deployment_id** (String, Optional) A filter to search by deployment ID.
- **environments** (List of String, Optional) A filter to search by a list of environment IDs.
- **health_statuses** (List of String, Optional) A filter to search by a list of health statuses of resources. Valid health statuses are `HasWarnings`, `Healthy`, `Unavailable`, `Unhealthy`, or `Unknown`.
- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **is_disabled** (Boolean, Optional) A filter to search by the disabled status of a resource.
- **name** (String, Optional) A filter to search by name.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **roles** (List of String, Optional) A filter to search by a list of role IDs.
- **shell_names** (List of String, Optional) A list of shell names to match in the query and/or search
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.
- **tenant_tags** (List of String, Optional) A filter to search by a list of tenant tags.
- **tenants** (List of String, Optional) A filter to search by a list of tenant IDs.
- **thumbprint** (String, Optional) The thumbprint of the deployment target to match in the query and/or search

### Read-only

- **azure_service_fabric_cluster_deployment_targets** (Block List) A list of Azure service fabric cluster deployment targets that match the filter(s). (see [below for nested schema](#nestedblock--azure_service_fabric_cluster_deployment_targets))
- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.

<a id="nestedblock--azure_service_fabric_cluster_deployment_targets"></a>
### Nested Schema for `azure_service_fabric_cluster_deployment_targets`

Read-only:

- **aad_client_credential_secret** (String, Read-only)
- **aad_credential_type** (String, Read-only)
- **aad_user_credential_password** (String, Read-only)
- **aad_user_credential_username** (String, Read-only)
- **certificate_store_location** (String, Read-only)
- **certificate_store_name** (String, Read-only)
- **client_certificate_variable** (String, Read-only)
- **connection_endpoint** (String, Read-only)
- **environments** (List of String, Read-only) A list of environment IDs associated with this resource.
- **has_latest_calamari** (Boolean, Read-only)
- **health_status** (String, Read-only) Represents the health status of this deployment target. Valid health statuses are `HasWarnings`, `Healthy`, `Unavailable`, `Unhealthy`, or `Unknown`.
- **id** (String, Read-only) The unique ID for this resource.
- **is_disabled** (Boolean, Read-only)
- **is_in_process** (Boolean, Read-only)
- **machine_policy_id** (String, Read-only)
- **name** (String, Read-only) The name of this resource.
- **operating_system** (String, Read-only)
- **roles** (List of String, Read-only)
- **security_mode** (String, Read-only)
- **server_certificate_thumbprint** (String, Read-only)
- **shell_name** (String, Read-only)
- **shell_version** (String, Read-only)
- **space_id** (String, Read-only) The space ID associated with this resource.
- **status** (String, Read-only) The status of this resource. Valid statuses are `CalamariNeedsUpgrade`, `Disabled`, `NeedsUpgrade`, `Offline`, `Online`, or `Unknown`.
- **status_summary** (String, Read-only) A summary elaborating on the status of this resource.
- **tenant_tags** (List of String, Read-only) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Read-only) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Read-only) A list of tenant IDs associated with this resource.
- **thumbprint** (String, Read-only)
- **uri** (String, Read-only)


