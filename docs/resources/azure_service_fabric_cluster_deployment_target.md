---
page_title: "octopusdeploy_azure_service_fabric_cluster_deployment_target Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages Azure service fabric cluster deployment targets in Octopus Deploy.
---

# Resource `octopusdeploy_azure_service_fabric_cluster_deployment_target`

This resource manages Azure service fabric cluster deployment targets in Octopus Deploy.



## Schema

### Required

- **connection_endpoint** (String, Required)
- **name** (String, Required) The name of this resource.
- **roles** (List of String, Required)

### Optional

- **aad_client_credential_secret** (String, Optional)
- **aad_credential_type** (String, Optional)
- **aad_user_credential_password** (String, Optional)
- **aad_user_credential_username** (String, Optional)
- **certificate_store_location** (String, Optional)
- **certificate_store_name** (String, Optional)
- **client_certificate_variable** (String, Optional)
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **health_status** (String, Optional) Represents the health status of this deployment target. Valid health statuses are `HasWarnings`, `Healthy`, `Unavailable`, `Unhealthy`, or `Unknown`.
- **id** (String, Optional) The unique identifier for this resource.
- **is_disabled** (Boolean, Optional)
- **machine_policy_id** (String, Optional)
- **operating_system** (String, Optional)
- **security_mode** (String, Optional)
- **server_certificate_thumbprint** (String, Optional)
- **shell_name** (String, Optional)
- **shell_version** (String, Optional)
- **space_id** (String, Optional) The space identifier associated with this resource.
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


