---
page_title: "octopusdeploy_kubernetes_cluster_deployment_target Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages Kubernets cluster deployment targets in Octopus Deploy.
---

# Resource `octopusdeploy_kubernetes_cluster_deployment_target`

This resource manages Kubernets cluster deployment targets in Octopus Deploy.



## Schema

### Required

- **cluster_url** (String, Required)
- **environments** (List of String, Required) A list of environment IDs associated with this resource.
- **name** (String, Required) The name of this resource.
- **roles** (List of String, Required)

### Optional

- **authentication** (Block List, Max: 1) (see [below for nested schema](#nestedblock--authentication))
- **aws_account_authentication** (Block List, Max: 1) (see [below for nested schema](#nestedblock--aws_account_authentication))
- **azure_service_principal_authentication** (Block List, Max: 1) (see [below for nested schema](#nestedblock--azure_service_principal_authentication))
- **certificate_authentication** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--certificate_authentication))
- **cluster_certificate** (String, Optional)
- **container** (Block List) (see [below for nested schema](#nestedblock--container))
- **default_worker_pool_id** (String, Optional)
- **endpoint** (Block List) (see [below for nested schema](#nestedblock--endpoint))
- **health_status** (String, Optional) Represents the health status of this deployment target. Valid health statuses are `HasWarnings`, `Healthy`, `Unavailable`, `Unhealthy`, or `Unknown`.
- **id** (String, Optional) The unique ID for this resource.
- **is_disabled** (Boolean, Optional)
- **machine_policy_id** (String, Optional)
- **namespace** (String, Optional)
- **operating_system** (String, Optional)
- **proxy_id** (String, Optional)
- **running_in_container** (Boolean, Optional)
- **shell_name** (String, Optional)
- **shell_version** (String, Optional)
- **skip_tls_verification** (Boolean, Optional)
- **space_id** (String, Optional) The space ID associated with this resource.
- **status** (String, Optional) The status of this resource. Valid statuses are `CalamariNeedsUpgrade`, `Disabled`, `NeedsUpgrade`, `Offline`, `Online`, or `Unknown`.
- **status_summary** (String, Optional) A summary elaborating on the status of this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.
- **thumbprint** (String, Optional)
- **uri** (String, Optional)

### Read-only

- **has_latest_calamari** (Boolean, Read-only)
- **is_in_process** (Boolean, Read-only)

<a id="nestedblock--authentication"></a>
### Nested Schema for `authentication`

Optional:

- **account_id** (String, Optional)


<a id="nestedblock--aws_account_authentication"></a>
### Nested Schema for `aws_account_authentication`

Required:

- **account_id** (String, Required)
- **cluster_name** (String, Required)

Optional:

- **assume_role** (Boolean, Optional)
- **assume_role_external_id** (String, Optional)
- **assume_role_session_duration** (Number, Optional)
- **assumed_role_arn** (String, Optional)
- **assumed_role_session** (String, Optional)
- **use_instance_role** (Boolean, Optional)


<a id="nestedblock--azure_service_principal_authentication"></a>
### Nested Schema for `azure_service_principal_authentication`

Required:

- **account_id** (String, Required)
- **cluster_name** (String, Required)
- **cluster_resource_group** (String, Required)


<a id="nestedblock--certificate_authentication"></a>
### Nested Schema for `certificate_authentication`

Optional:

- **client_certificate** (String, Optional)


<a id="nestedblock--container"></a>
### Nested Schema for `container`

Optional:

- **feed_id** (String, Optional)
- **image** (String, Optional)


<a id="nestedblock--endpoint"></a>
### Nested Schema for `endpoint`

Required:

- **communication_style** (String, Required)

Optional:

- **aad_client_credential_secret** (String, Optional)
- **aad_credential_type** (String, Optional)
- **aad_user_credential_username** (String, Optional)
- **account_id** (String, Optional)
- **applications_directory** (String, Optional)
- **authentication** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--endpoint--authentication))
- **certificate_signature_algorithm** (String, Optional)
- **certificate_store_location** (String, Optional)
- **certificate_store_name** (String, Optional)
- **client_certificate_variable** (String, Optional)
- **cloud_service_name** (String, Optional)
- **cluster_certificate** (String, Optional)
- **cluster_url** (String, Optional)
- **connection_endpoint** (String, Optional)
- **container** (Block List) (see [below for nested schema](#nestedblock--endpoint--container))
- **default_worker_pool_id** (String, Optional)
- **destination** (Block List) (see [below for nested schema](#nestedblock--endpoint--destination))
- **dot_net_core_platform** (String, Optional)
- **fingerprint** (String, Optional)
- **host** (String, Optional)
- **id** (String, Optional) The unique ID for this resource.
- **namespace** (String, Optional)
- **port** (Number, Optional)
- **proxy_id** (String, Optional)
- **resource_group_name** (String, Optional)
- **running_in_container** (Boolean, Optional)
- **security_mode** (String, Optional)
- **server_certificate_thumbprint** (String, Optional)
- **skip_tls_verification** (Boolean, Optional)
- **slot** (String, Optional)
- **storage_account_name** (String, Optional)
- **swap_if_possible** (Boolean, Optional)
- **tentacle_version_details** (Block List) (see [below for nested schema](#nestedblock--endpoint--tentacle_version_details))
- **thumbprint** (String, Optional)
- **uri** (String, Optional)
- **use_current_instance_count** (Boolean, Optional)
- **web_app_name** (String, Optional)
- **web_app_slot_name** (String, Optional)
- **working_directory** (String, Optional)

<a id="nestedblock--endpoint--authentication"></a>
### Nested Schema for `endpoint.authentication`

Optional:

- **account_id** (String, Optional)
- **admin_login** (String, Optional)
- **assume_role** (Boolean, Optional)
- **assume_role_external_id** (String, Optional)
- **assume_role_session_duration** (Number, Optional)
- **assumed_role_arn** (String, Optional)
- **assumed_role_session** (String, Optional)
- **authentication_type** (String, Optional)
- **client_certificate** (String, Optional)
- **cluster_name** (String, Optional)
- **cluster_resource_group** (String, Optional)
- **use_instance_role** (Boolean, Optional)


<a id="nestedblock--endpoint--container"></a>
### Nested Schema for `endpoint.container`

Optional:

- **feed_id** (String, Optional)
- **image** (String, Optional)


<a id="nestedblock--endpoint--destination"></a>
### Nested Schema for `endpoint.destination`

Optional:

- **destination_type** (String, Optional)
- **drop_folder_path** (String, Optional)


<a id="nestedblock--endpoint--tentacle_version_details"></a>
### Nested Schema for `endpoint.tentacle_version_details`

Optional:

- **upgrade_locked** (Boolean, Optional)
- **upgrade_required** (Boolean, Optional)
- **upgrade_suggested** (Boolean, Optional)
- **version** (String, Optional)


