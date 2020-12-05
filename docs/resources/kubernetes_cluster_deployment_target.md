---
page_title: "octopusdeploy_kubernetes_cluster_deployment_target Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  
---

# Resource `octopusdeploy_kubernetes_cluster_deployment_target`





## Schema

### Required

- **cluster_url** (String, Required)
- **name** (String, Required) The name of this resource.
- **roles** (List of String, Required)

### Optional

- **aws_account_authentication** (Block List, Max: 1) (see [below for nested schema](#nestedblock--aws_account_authentication))
- **azure_service_principal_authentication** (Block List, Max: 1) (see [below for nested schema](#nestedblock--azure_service_principal_authentication))
- **certificate_authentication** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--certificate_authentication))
- **cluster_certificate** (String, Optional)
- **container** (Block List) (see [below for nested schema](#nestedblock--container))
- **default_worker_pool_id** (String, Optional)
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **health_status** (String, Optional) Represents the health status of this deployment target. Valid health statuses are `HasWarnings`, `Healthy`, `Unavailable`, `Unhealthy`, or `Unknown`.
- **id** (String, Optional) The unique identifier for this resource.
- **is_disabled** (Boolean, Optional)
- **machine_policy_id** (String, Optional)
- **namespace** (String, Optional)
- **operating_system** (String, Optional)
- **proxy_id** (String, Optional)
- **running_in_container** (Boolean, Optional)
- **shell_name** (String, Optional)
- **shell_version** (String, Optional)
- **skip_tls_verification** (Boolean, Optional)
- **space_id** (String, Optional) The space identifier associated with this resource.
- **status** (String, Optional) The status of this resource. Valid statuses are `CalamariNeedsUpgrade`, `Disabled`, `NeedsUpgrade`, `Offline`, `Online`, or `Unknown`.
- **status_summary** (String, Optional) A summary elaborating on the status of this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.
- **thumbprint** (String, Optional)
- **token_authentication** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--token_authentication))
- **uri** (String, Optional)
- **username_password_authentication** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--username_password_authentication))

### Read-only

- **has_latest_calamari** (Boolean, Read-only)
- **is_in_process** (Boolean, Read-only)

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


<a id="nestedblock--token_authentication"></a>
### Nested Schema for `token_authentication`

Optional:

- **account_id** (String, Optional)


<a id="nestedblock--username_password_authentication"></a>
### Nested Schema for `username_password_authentication`

Optional:

- **account_id** (String, Optional)

