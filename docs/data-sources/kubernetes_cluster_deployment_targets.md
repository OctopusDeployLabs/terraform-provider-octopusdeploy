---
page_title: "octopusdeploy_kubernetes_cluster_deployment_targets Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing Kubernetes cluster deployment targets.
---

# Data Source `octopusdeploy_kubernetes_cluster_deployment_targets`

Provides information about existing Kubernetes cluster deployment targets.



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

- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.
- **kubernetes_cluster_deployment_target** (Block List) A list of Kubernetes cluster deployment targets that match the filter(s). (see [below for nested schema](#nestedblock--kubernetes_cluster_deployment_target))

<a id="nestedblock--kubernetes_cluster_deployment_target"></a>
### Nested Schema for `kubernetes_cluster_deployment_target`

Read-only:

- **aws_account_authentication** (List of Object, Read-only) (see [below for nested schema](#nestedatt--kubernetes_cluster_deployment_target--aws_account_authentication))
- **azure_service_principal_authentication** (List of Object, Read-only) (see [below for nested schema](#nestedatt--kubernetes_cluster_deployment_target--azure_service_principal_authentication))
- **certificate_authentication** (Set of Object, Read-only) (see [below for nested schema](#nestedatt--kubernetes_cluster_deployment_target--certificate_authentication))
- **cluster_certificate** (String, Read-only)
- **cluster_url** (String, Read-only)
- **container** (List of Object, Read-only) (see [below for nested schema](#nestedatt--kubernetes_cluster_deployment_target--container))
- **default_worker_pool_id** (String, Read-only)
- **environments** (List of String, Read-only) A list of environment IDs associated with this resource.
- **has_latest_calamari** (Boolean, Read-only)
- **health_status** (String, Read-only) Represents the health status of this deployment target. Valid health statuses are `HasWarnings`, `Healthy`, `Unavailable`, `Unhealthy`, or `Unknown`.
- **id** (String, Read-only) The unique identifier for this resource.
- **is_disabled** (Boolean, Read-only)
- **is_in_process** (Boolean, Read-only)
- **machine_policy_id** (String, Read-only)
- **name** (String, Read-only) The name of this resource.
- **namespace** (String, Read-only)
- **operating_system** (String, Read-only)
- **proxy_id** (String, Read-only)
- **roles** (List of String, Read-only)
- **running_in_container** (Boolean, Read-only)
- **shell_name** (String, Read-only)
- **shell_version** (String, Read-only)
- **skip_tls_verification** (Boolean, Read-only)
- **space_id** (String, Read-only) The space identifier associated with this resource.
- **status** (String, Read-only) The status of this resource. Valid statuses are `CalamariNeedsUpgrade`, `Disabled`, `NeedsUpgrade`, `Offline`, `Online`, or `Unknown`.
- **status_summary** (String, Read-only) A summary elaborating on the status of this resource.
- **tenant_tags** (List of String, Read-only) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Read-only) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Read-only) A list of tenant IDs associated with this resource.
- **thumbprint** (String, Read-only)
- **token_authentication** (Set of Object, Read-only) (see [below for nested schema](#nestedatt--kubernetes_cluster_deployment_target--token_authentication))
- **uri** (String, Read-only)
- **username_password_authentication** (Set of Object, Read-only) (see [below for nested schema](#nestedatt--kubernetes_cluster_deployment_target--username_password_authentication))

<a id="nestedatt--kubernetes_cluster_deployment_target--aws_account_authentication"></a>
### Nested Schema for `kubernetes_cluster_deployment_target.aws_account_authentication`

- **account_id** (String)
- **assume_role** (Boolean)
- **assume_role_external_id** (String)
- **assume_role_session_duration** (Number)
- **assumed_role_arn** (String)
- **assumed_role_session** (String)
- **cluster_name** (String)
- **use_instance_role** (Boolean)


<a id="nestedatt--kubernetes_cluster_deployment_target--azure_service_principal_authentication"></a>
### Nested Schema for `kubernetes_cluster_deployment_target.azure_service_principal_authentication`

- **account_id** (String)
- **cluster_name** (String)
- **cluster_resource_group** (String)


<a id="nestedatt--kubernetes_cluster_deployment_target--certificate_authentication"></a>
### Nested Schema for `kubernetes_cluster_deployment_target.certificate_authentication`

- **client_certificate** (String)


<a id="nestedatt--kubernetes_cluster_deployment_target--container"></a>
### Nested Schema for `kubernetes_cluster_deployment_target.container`

- **feed_id** (String)
- **image** (String)


<a id="nestedatt--kubernetes_cluster_deployment_target--token_authentication"></a>
### Nested Schema for `kubernetes_cluster_deployment_target.token_authentication`

- **account_id** (String)


<a id="nestedatt--kubernetes_cluster_deployment_target--username_password_authentication"></a>
### Nested Schema for `kubernetes_cluster_deployment_target.username_password_authentication`

- **account_id** (String)

