---
page_title: "octopusdeploy_azure_web_app_deployment_targets Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing Azure web app deployment targets.
---

# Data Source `octopusdeploy_azure_web_app_deployment_targets`

Provides information about existing Azure web app deployment targets.

## Example Usage

```terraform
data "octopusdeploy_azure_web_app_deployment_targets" "example" {
  health_statuses = ["Healthy", "Unavailable"]
  ids             = ["Machines-123", "Machines-321"]
  partial_name    = "Defau"
  skip            = 5
  take            = 100
}
```

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

- **azure_web_app_deployment_targets** (Block List) A list of Azure web app deployment targets that match the filter(s). (see [below for nested schema](#nestedblock--azure_web_app_deployment_targets))
- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.

<a id="nestedblock--azure_web_app_deployment_targets"></a>
### Nested Schema for `azure_web_app_deployment_targets`

Read-only:

- **account_id** (String, Read-only)
- **endpoint** (List of Object, Read-only) (see [below for nested schema](#nestedatt--azure_web_app_deployment_targets--endpoint))
- **environments** (List of String, Read-only) A list of environment IDs associated with this resource.
- **has_latest_calamari** (Boolean, Read-only)
- **health_status** (String, Read-only) Represents the health status of this deployment target. Valid health statuses are `HasWarnings`, `Healthy`, `Unavailable`, `Unhealthy`, or `Unknown`.
- **id** (String, Read-only) The unique ID for this resource.
- **is_disabled** (Boolean, Read-only)
- **is_in_process** (Boolean, Read-only)
- **machine_policy_id** (String, Read-only)
- **name** (String, Read-only) The name of this resource.
- **operating_system** (String, Read-only)
- **resource_group_name** (String, Read-only)
- **roles** (List of String, Read-only)
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
- **web_app_name** (String, Read-only)
- **web_app_slot_name** (String, Read-only)

<a id="nestedatt--azure_web_app_deployment_targets--endpoint"></a>
### Nested Schema for `azure_web_app_deployment_targets.endpoint`

- **aad_client_credential_secret** (String)
- **aad_credential_type** (String)
- **aad_user_credential_username** (String)
- **account_id** (String)
- **applications_directory** (String)
- **authentication** (Set of Object) (see [below for nested schema](#nestedobjatt--azure_web_app_deployment_targets--endpoint--authentication))
- **certificate_signature_algorithm** (String)
- **certificate_store_location** (String)
- **certificate_store_name** (String)
- **client_certificate_variable** (String)
- **cloud_service_name** (String)
- **cluster_certificate** (String)
- **cluster_url** (String)
- **communication_style** (String)
- **connection_endpoint** (String)
- **container** (List of Object) (see [below for nested schema](#nestedobjatt--azure_web_app_deployment_targets--endpoint--container))
- **default_worker_pool_id** (String)
- **destination** (List of Object) (see [below for nested schema](#nestedobjatt--azure_web_app_deployment_targets--endpoint--destination))
- **dot_net_core_platform** (String)
- **fingerprint** (String)
- **host** (String)
- **id** (String)
- **namespace** (String)
- **port** (Number)
- **proxy_id** (String)
- **resource_group_name** (String)
- **running_in_container** (Boolean)
- **security_mode** (String)
- **server_certificate_thumbprint** (String)
- **skip_tls_verification** (Boolean)
- **slot** (String)
- **storage_account_name** (String)
- **swap_if_possible** (Boolean)
- **tentacle_version_details** (List of Object) (see [below for nested schema](#nestedobjatt--azure_web_app_deployment_targets--endpoint--tentacle_version_details))
- **thumbprint** (String)
- **uri** (String)
- **use_current_instance_count** (Boolean)
- **web_app_name** (String)
- **web_app_slot_name** (String)
- **working_directory** (String)

<a id="nestedobjatt--azure_web_app_deployment_targets--endpoint--authentication"></a>
### Nested Schema for `azure_web_app_deployment_targets.endpoint.authentication`

- **account_id** (String)
- **admin_login** (String)
- **assume_role** (Boolean)
- **assume_role_external_id** (String)
- **assume_role_session_duration** (Number)
- **assumed_role_arn** (String)
- **assumed_role_session** (String)
- **authentication_type** (String)
- **client_certificate** (String)
- **cluster_name** (String)
- **cluster_resource_group** (String)
- **use_instance_role** (Boolean)


<a id="nestedobjatt--azure_web_app_deployment_targets--endpoint--container"></a>
### Nested Schema for `azure_web_app_deployment_targets.endpoint.container`

- **feed_id** (String)
- **image** (String)


<a id="nestedobjatt--azure_web_app_deployment_targets--endpoint--destination"></a>
### Nested Schema for `azure_web_app_deployment_targets.endpoint.destination`

- **destination_type** (String)
- **drop_folder_path** (String)


<a id="nestedobjatt--azure_web_app_deployment_targets--endpoint--tentacle_version_details"></a>
### Nested Schema for `azure_web_app_deployment_targets.endpoint.tentacle_version_details`

- **upgrade_locked** (Boolean)
- **upgrade_required** (Boolean)
- **upgrade_suggested** (Boolean)
- **version** (String)


